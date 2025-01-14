package main

import (
	"os"

	"github.com/spf13/cobra"
)

func PatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "patch",
		Aliases: []string{"p", "pa"},
		Short:   "Bump the patch version",
		RunE:    Patch,
	}
}

func Patch(cmd *cobra.Command, args []string) error {
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

	newVersion := version.BumpPatch()
	if *Prefix != "" {
		newVersion.Prefix = Prefix
	}
	Info("%s -> %s\n", version.String(), newVersion.String())

	if *DryRun {
		Info("dry run, will not create tag\n")
		return nil
	}

	return repo.CreateAndPushTag(newVersion.String())
}
