package main

import "github.com/spf13/cobra"

var BumpVersion = "dev"

func BumpVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of bump",
		Run: func(cmd *cobra.Command, args []string) {
			Info("bump %s\n", BumpVersion)
		},
	}
}
