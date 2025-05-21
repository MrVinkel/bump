package main

import (
	"os"

	"github.com/mrvinkel/bump/cmd/bump/internal"
	"github.com/spf13/cobra"
)

var (
	patchCmd = &cobra.Command{
		Use:     "patch",
		Aliases: []string{"p", "pa"},
		Short:   "Bump the patch version",
		RunE:    internal.Bump(internal.BumpPatch),
	}
	minorCmd = &cobra.Command{
		Use:     "minor",
		Short:   "Bump the minor version",
		Aliases: []string{"m", "mi"},
		RunE:    internal.Bump(internal.BumpMinor),
	}
	majorCmd = &cobra.Command{
		Use:     "major",
		Short:   "Bump the major version",
		Aliases: []string{"M", "ma"},
		RunE:    internal.Bump(internal.BumpMajor),
	}
	preReleaseCmd = &cobra.Command{
		Use:     "prerelease",
		Short:   "Bump the pre-release version",
		Aliases: []string{"pre"},
		RunE:    internal.BumpE(internal.BumpPreRelease),
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of bump",
		Run: func(cmd *cobra.Command, args []string) {
			internal.Info("bump %s\n", internal.BumpVersion)
		},
	}
)

func main() {
	root := &cobra.Command{
		Use:           "bump",
		Short:         "Bump those versions!",
		Long:          `Bump those versions! Utility for bumping and pushing git tags`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          internal.Bump(internal.BumpPatch),
	}

	internal.DebugFlag = root.PersistentFlags().BoolP("debug", "d", false, "Debug mode")
	internal.QuietFlag = root.PersistentFlags().BoolP("quiet", "q", false, "Quiet - only output errors")
	internal.DryRun = root.PersistentFlags().BoolP("dry-run", "x", false, "Do not create tags, only print what would be done")
	internal.NoVerify = root.PersistentFlags().BoolP("no-verify", "n", false, "Do not check repository status before creating tags")
	internal.NoFetch = root.PersistentFlags().BoolP("no-fetch", "f", false, "Do not fetch before verifying repository status")
	internal.NoCommit = root.PersistentFlags().BoolP("no-commit", "c", false, "Do not commit changes to the repository")
	internal.SkipPreHook = root.PersistentFlags().BoolP("skip-pre-hook", "s", false, "Skip any configured pre-hook")
	internal.Prefix = root.PersistentFlags().StringP("prefix", "p", "", "Prefix for the version tag")
	internal.Build = root.PersistentFlags().String("build", "", "Build metadata to prepend to the version tag")
	internal.Alpha = root.PersistentFlags().BoolP("alpha", "a", false, "Bump the pre-release version to alpha.1")
	internal.Beta = root.PersistentFlags().BoolP("beta", "b", false, "Bump the pre-release version to beta.1")
	internal.RC = root.PersistentFlags().BoolP("rc", "r", false, "Bump the pre-release version to rc.1")

	root.AddCommand(versionCmd)
	root.AddCommand(patchCmd)
	root.AddCommand(minorCmd)
	root.AddCommand(majorCmd)
	root.AddCommand(preReleaseCmd)

	if err := root.Execute(); err != nil {
		internal.Error("%v\n", err)
		os.Exit(1)
	}
}
