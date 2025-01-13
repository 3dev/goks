package main

import (
	"fmt"
	"github.com/3dev/goKeyStore"
	"github.com/spf13/cobra"
	"path/filepath"
)

func buildStatsCommand(rootCmd *cobra.Command) {

	var passkey string
	var filename string

	var createCmd = &cobra.Command{
		Use:   "stats",
		Short: "provides the statistic of the go keystore file",
		Long:  "displays statistics regarding the go keystore file",
		Run: func(cmd *cobra.Command, args []string) {

			ext := filepath.Ext(filename)

			// If there's no extension, add ".gks"
			if ext == "" {
				filename += ".gks"
			}

			_, err := goKeyStore.Open(filename, passkey)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("go keystore file '%s' statistics \n", filename)
		},
	}

	// Add flags to the config command
	createCmd.Flags().StringVar(&passkey, "passkey", "", "keystore access passkey")
	createCmd.Flags().StringVar(&filename, "file", "", "filename for the  keystore")

	_ = createCmd.MarkFlagRequired("passkey")
	_ = createCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(createCmd)
}
