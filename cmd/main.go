package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "goks",
		Short: "A tool create, display and manage encrypted content",
		Long:  "This tool provides commands to generate/inspect/modify a golang keystore basically storing key=value.",
	}

	buildVersionCommand(rootCmd)
	buildCreateCommand(rootCmd)
	buildStatsCommand(rootCmd)
	buildHexCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
