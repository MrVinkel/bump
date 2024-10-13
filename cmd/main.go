package main

import (
	"bump/cmd/bump"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "bump",
		Short:         "Bump those versions!",
		Long:          `Utility for bumping and pushing git tags`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	bump.DebugFlag = root.PersistentFlags().BoolP("debug", "d", false, "Debug mode")
	bump.QuietFlag = root.PersistentFlags().BoolP("quiet", "q", false, "Quiet - only output errors")
	bump.DryRun = root.PersistentFlags().BoolP("dry-run", "n", false, "Dry run mode")

	root.AddCommand(bump.PatchCmd())
	root.AddCommand(bump.MinorCmd())
	root.AddCommand(bump.MajorCmd())

	if err := root.Execute(); err != nil {
		bump.Error("%v\n", err)
		os.Exit(1)
	}
}
