package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

func buildHexCommand(rootCmd *cobra.Command) {

	var passkey string
	var filename string
	var key string

	var createCmd = &cobra.Command{
		Use:   "hex",
		Short: "hex of a key",
		Long:  "displays the hexdump of a key within the go keystore file",
		Run: func(cmd *cobra.Command, args []string) {

			ks, err := openKeyStore(filename, passkey)
			if err != nil {
				fmt.Printf("unable to open keystore: %v\n", err)
				return
			}

			if !slices.Contains(ks.Keys(), key) {
				fmt.Printf("key '%s' not found in keystore\n", key)
				return
			}

			kInfo, err := ks.KeyInfo(key)
			if err != nil {
				fmt.Printf("unable to get key info: %v\n", err)
			}

			fmt.Println()
			fmt.Printf("key info:\n")
			fmt.Printf("  index available:\t%v\n", kInfo.Available < 0)
			fmt.Printf("  key name:\t\t%s\n", key)
			fmt.Printf("  data length:\t\t%d bytes \t[%.2f KB]\n", binary.BigEndian.Uint32(kInfo.DataLength[:]), float64(binary.BigEndian.Uint32(kInfo.DataLength[:]))/1024.0)
			fmt.Printf("  allocated space:\t%d bytes \t[%.2f KB]\n", binary.BigEndian.Uint32(kInfo.AllocatedLength[:]), float64(binary.BigEndian.Uint32(kInfo.AllocatedLength[:]))/1024.0)
			fmt.Printf("  file position:\t%d bytes \t[%.2f KB]\n", binary.BigEndian.Uint32(kInfo.Location[:]), float64(binary.BigEndian.Uint32(kInfo.Location[:]))/1024.0)

			data, err := ks.Get(key)
			if err != nil {
				fmt.Printf("unable to get key: %v\n", err)
				return
			}

			fmt.Println()
			fmt.Printf("data:\n")
			fmt.Printf("%s\n", hex.Dump(data))

		},
	}

	// Add flags to the config command
	createCmd.Flags().StringVar(&passkey, "passkey", "", "keystore access passkey")
	createCmd.Flags().StringVar(&filename, "file", "", "filename for the  keystore")
	createCmd.Flags().StringVar(&key, "key", "", "get details for the key")

	_ = createCmd.MarkFlagRequired("passkey")
	_ = createCmd.MarkFlagRequired("file")
	_ = createCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(createCmd)
}
