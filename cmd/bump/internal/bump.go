package internal

import (
	"errors"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var (
	BumpVersion = "dev"

	DebugFlag   *bool
	QuietFlag   *bool
	DryRun      *bool
	NoVerify    *bool
	NoFetch     *bool
	NoCommit    *bool
	SkipPreHook *bool
	Prefix      *string
	Build       *string
	Alpha       *bool
	Beta        *bool
	RC          *bool
)

func Bump(fn func(*Version) *Version) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return bump(func(v *Version) (*Version, error) {
			return fn(v), nil
		})
	}
}

func BumpE(fn func(*Version) (*Version, error)) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return bump(fn)
	}
}

func bump(fn func(*Version) (*Version, error)) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := NewRepo(cwd)
	if err != nil {
		return err
	}

	if !*NoVerify {
		err = checkRepositoryStatus(repo)
		if err != nil {
			return err
		}
	}

	repoDir, err := repo.GetDir()
	if err != nil {
		return err
	}

	config, err := ReadConfig(os.DirFS(repoDir))
	if err != nil {
		return err
	}
	useConfig(config)

	previousVersion, err := getLatestVersion(repo)
	if err != nil {
		return err
	}

	newVersion, err := fn(previousVersion)
	if err != nil {
		return err
	}

	// reuse the prefix unless set
	if Prefix != nil && *Prefix != "" {
		newVersion.Prefix = Prefix
	}

	// never reuse build metadata
	if Build != nil && *Build != "" {
		newVersion.Build = Build
	} else {
		newVersion.Build = nil
	}

	preCount := 0
	if Alpha != nil && *Alpha {
		Debug("bumping to alpha\n")
		newVersion.Alpha()
		preCount++
	}
	if Beta != nil && *Beta {
		Debug("bumping to beta\n")
		newVersion.Beta()
		preCount++
	}
	if RC != nil && *RC {
		Debug("bumping to rc\n")
		newVersion.RC()
		preCount++
	}

	if preCount > 1 {
		return errors.New("only one of --alpha, --beta, --rc can be specified")
	}

	err = runPreHook(config, newVersion, previousVersion)
	if err != nil {
		Debug("error: %v\n", err)
		return errors.New("pre-hook failed")
	}

	err = commitChanges(config, repo, newVersion, previousVersion)
	if err != nil {
		return err
	}

	Info("tag: %s -> %s\n", previousVersion.String(), newVersion.String())

	if *DryRun {
		Info("dry run, will not create tag\n")
		return nil
	}

	return repo.TagAndPush(newVersion.String())
}

func useConfig(config *Config) {
	if config == nil {
		return
	}

	if config.Verify != nil {
		*NoVerify = !*config.Verify
	}
	if config.Fetch != nil {
		*NoFetch = !*config.Fetch
	}
	if config.Commit != nil {
		*NoCommit = !*config.Commit
	}
	if config.Prefix != nil {
		*Prefix = *config.Prefix
	}
}

func getLatestVersion(repo *Repo) (*Version, error) {
	tags, err := repo.GetTags()
	if err != nil {
		return nil, err
	}

	Debug("tags: %v\n", SliceString(tags))

	versions := make([]Version, 0)
	for _, t := range tags {
		v, err := ParseVersion(t)
		if err != nil {
			Error("invalid tag: %s\n", t)
			Debug("error: %v\n", err)
			continue
		}
		versions = append(versions, *v)
	}

	Debug("parsed versions: %s\n", VersionSliceString(versions))

	sort.Slice(versions, func(i, j int) bool {
		return Compare(versions[i], versions[j]) > 0
	})

	Debug("sorted versions: %v\n", VersionSliceString(versions))

	version := NewVersion(nil, 0, 0, 0, []string{}, nil)
	if len(versions) > 0 {
		version = &versions[0]
	}
	return version, nil
}

func checkRepositoryStatus(repo *Repo) error {
	hasChanages, err := repo.HasChanges()
	if err != nil {
		return err
	}
	if hasChanages {
		return errors.New("uncommitted changes")
	}

	if !*NoFetch {
		Debug("fetching repository\n")
		if err = repo.Fetch(); err != nil {
			return err
		}
	}

	synced, err := repo.IsSynced()
	if err != nil {
		return err
	}
	if !synced {
		return errors.New("unpushed changes")
	}
	return nil
}

func runPreHook(config *Config, newVersion, previousVersion *Version) error {
	if config == nil || len(config.PreHook) == 0 {
		return nil
	}
	if *SkipPreHook {
		Info("skipping pre hook\n")
		return nil
	}

	Debug("running pre hook\n")
	env := map[string]string{
		"VERSION":          newVersion.String(),
		"PREVIOUS_VERSION": previousVersion.String(),
	}
	Info("running pre-hook\n")
	return Run(*config.Shell, config.PreHook, os.Stdout, env)
}

func commitChanges(config *Config, repo *Repo, newVersion, previousVersion *Version) error {
	if config == nil {
		return nil
	}
	if *NoCommit {
		return nil
	}
	err := checkRepositoryStatus(repo)
	if err == nil {
		// no changes, nothing to commit
		return nil
	}

	vars := map[string]string{
		"VERSION":          newVersion.String(),
		"PREVIOUS_VERSION": previousVersion.String(),
	}
	message := os.Expand(*config.Message, func(s string) string {
		return vars[s]
	})

	Info("commit: %s\n", message)
	if *DryRun {
		Info("dry run, will not commit and push changes\n")
		return nil
	}
	return repo.CommitAndPush(message)
}
