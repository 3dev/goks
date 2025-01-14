package main

import (
	"fmt"
	"github.com/3dev/goKeyStore"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func buildCreateCommand(rootCmd *cobra.Command) {

	var passkey string
	var filename string
	var overwrite bool

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "creates a new go keystore file",
		Long:  "crates a new go keystore file",
		Run: func(cmd *cobra.Command, args []string) {

			ext := filepath.Ext(filename)

			// If there's no extension, add ".gks"
			if ext == "" {
				filename += ".goks"
			}

			_, err := os.Stat(filename)
			if !os.IsNotExist(err) && !overwrite {
				fmt.Printf("'%s' already exists, use --overwrite flag to overwrite the file\n", filename)
				return
			}

			_, err = goKeyStore.New(filename, passkey)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Created go keystore file '%s' successfully\n", filename)
		},
	}

	// Add flags to the config command
	createCmd.Flags().StringVar(&passkey, "passkey", "", "keystore access passkey")
	createCmd.Flags().StringVar(&filename, "file", "", "filename for the  keystore")
	createCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite the keystore file")
	_ = createCmd.MarkFlagRequired("passkey")
	_ = createCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(createCmd)
}
