package bump

import (
	"os"

	"github.com/spf13/cobra"
)

func PatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "patch",
		Aliases: []string{"p"},
		Short:   "Bump the patch version",
		RunE:    patch,
	}
}

func patch(cmd *cobra.Command, args []string) error {
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

	newVersion := version.BumpPatch()
	Info("%s -> %s\n", version.String(), newVersion.String())

	if *DryRun {
		Info("dry run, not creating tags\n")
		return nil
	}

	return repo.CreateAndPushTag(newVersion.String())

}
