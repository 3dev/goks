package main

import (
	"github.com/3dev/goks"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func buildAddCommand(rootCmd *cobra.Command) {

	var passkey string
	var filename string
	var key string
	var value string

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "add a new entry into the go keystore file",
		Long:  "add a new entry into go keystore file",
		Run: func(cmd *cobra.Command, args []string) {

			ext := filepath.Ext(filename)

			// If there's no extension, add ".gks"
			if ext == "" {
				filename += ".goks"
			}

			ks, err := goks.Open(filename, passkey)
			if err != nil {
				color.Red("error opening go keystore file (%s): %v", filename, err)
				return
			}
			defer ks.Close()

			tokens := strings.Split(value, "=")
			if len(tokens) == 2 {
				if "file" != strings.ToLower(tokens[0]) {
					color.Red("error: when using the '=' separator for value it needs to be file={filename}")
					return
				}

				fileContent, err := os.ReadFile(tokens[1])
				if err != nil {
					color.Red("error reading file (%s): %v", tokens[1], err)
					return
				}

				err = ks.Put(key, fileContent)
				if err != nil {
					color.Red("error writing file (%s): %v", tokens[1], err)
					return
				}

				color.Green("added file (%s) into keystore", tokens[1])
				return
			}

			err = ks.Put(key, []byte(value))
			if err != nil {
				color.Red("error writing file (%s): %v", value, err)
			}

			color.Green("added value into keystore")

		},
	}

	// Add flags to the config command
	addCmd.Flags().StringVar(&passkey, "pass", "", "keystore access password")
	addCmd.Flags().StringVar(&filename, "file", "", "filename for the  keystore")
	addCmd.Flags().StringVar(&key, "key", "", "the item's key")
	addCmd.Flags().StringVar(&value, "value", "", "the item's value")

	_ = addCmd.MarkFlagRequired("pass")
	_ = addCmd.MarkFlagRequired("file")
	_ = addCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(addCmd)
}
