package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var (
	version string
)

func buildVersionCommand(rootCmd *cobra.Command) {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version",
		Long:  "Show the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("goks version is '%s'\n", strings.TrimSpace(version))
		},
	}

	rootCmd.AddCommand(versionCmd)
}
