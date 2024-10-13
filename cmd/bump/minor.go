package bump

import (
	"os"

	"github.com/spf13/cobra"
)

func MinorCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "minor",
		Aliases: []string{"m"},
		Short:   "Bump the minor version",
		RunE:    minor,
	}
}

func minor(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repo, err := NewRepo(cwd)
	if err != nil {
		return err
	}
	version, err := GetLatestVersion(repo)
	if err != nil {
		return err
	}

	newVersion := version.BumpMinor()
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
