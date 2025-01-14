package main

import (
	"os"

	"github.com/mrvinkel/bump/cmd/bump/internal"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "bump",
		Short:         "Bump those versions!",
		Long:          `Bump those versions! Utility for bumping and pushing git tags`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          internal.Patch,
	}

	internal.DebugFlag = root.PersistentFlags().BoolP("debug", "d", false, "Debug mode")
	internal.QuietFlag = root.PersistentFlags().BoolP("quiet", "q", false, "Quiet - only output errors")
	internal.DryRun = root.PersistentFlags().BoolP("dry-run", "x", false, "Do not create tags, only print what would be done")
	internal.NoVerify = root.PersistentFlags().BoolP("no-verify", "n", false, "Do not check repository status before creating tags")
	internal.NoFetch = root.PersistentFlags().BoolP("no-fetch", "f", false, "Do not fetch before verifying repository status")
	internal.Prefix = root.PersistentFlags().StringP("prefix", "p", "", "Prefix for the version tag")

	root.AddCommand(internal.BumpVersionCmd())
	root.AddCommand(internal.PatchCmd())
	root.AddCommand(internal.MinorCmd())
	root.AddCommand(internal.MajorCmd())

	if err := root.Execute(); err != nil {
		internal.Error("%v\n", err)
		os.Exit(1)
	}
}
