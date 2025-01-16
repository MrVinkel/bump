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
)

func Bump(fn func(*Version) *Version) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return bump(fn)
	}
}

func bump(fn func(*Version) *Version) error {
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

	newVersion := fn(previousVersion)
	if *Prefix != "" {
		newVersion.Prefix = Prefix
	}

	err = runPreHook(config, newVersion, previousVersion)
	if err != nil {
		Debug("error: %v\n", err)
		return errors.New("pre-hook failed")
	}

	err = commitChanges(repo, newVersion, previousVersion, *config.Message)
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

	version := NewVersion(0, 0, 0)
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

func commitChanges(repo *Repo, newVersion, previousVersion *Version, message string) error {
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
	message = os.Expand(message, func(s string) string {
		return vars[s]
	})

	Info("commit: %s\n", message)
	if *DryRun {
		Info("dry run, will not commit and push changes\n")
		return nil
	}
	return repo.CommitAndPush(message)
}
