package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "bump",
		Short:         "Bump those versions!",
		Long:          `Bump those versions! Utility for bumping and pushing git tags`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          Patch,
	}

	DebugFlag = root.PersistentFlags().BoolP("debug", "d", false, "Debug mode")
	QuietFlag = root.PersistentFlags().BoolP("quiet", "q", false, "Quiet - only output errors")
	DryRun = root.PersistentFlags().BoolP("dry-run", "x", false, "Do not create tags, only print what would be done")
	NoVerify = root.PersistentFlags().BoolP("no-verify", "n", false, "Do not check repository status before creating tags")
	NoFetch = root.PersistentFlags().BoolP("no-fetch", "f", false, "Do not fetch before verifying repository status")
	Prefix = root.PersistentFlags().StringP("prefix", "p", "", "Prefix for the version tag")

	root.AddCommand(BumpVersionCmd())
	root.AddCommand(PatchCmd())
	root.AddCommand(MinorCmd())
	root.AddCommand(MajorCmd())

	if err := root.Execute(); err != nil {
		Error("%v\n", err)
		os.Exit(1)
	}
}
