package internal

import (
	"os"

	"github.com/spf13/cobra"
)

func MajorCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "major",
		Aliases: []string{"ma"},
		Short:   "Bump the major version",
		RunE:    major,
	}
}

func major(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := NewRepo(cwd)
	if err != nil {
		return err
	}

	err = CheckRepositoryStatus(repo)
	if err != nil {
		return err
	}

	version, err := GetLatestVersion(repo)
	if err != nil {
		return err
	}

	newVersion := version.BumpMajor()
	if *Prefix != "" {
		newVersion.Prefix = Prefix
	}
	Info("%s -> %s\n", version.String(), newVersion.String())

	if *DryRun {
		Info("dry run, not creating tags\n")
		return nil
	}

	return repo.CreateAndPushTag(newVersion.String())
}
