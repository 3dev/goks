package main

import (
	"encoding/binary"
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

func buildStatsCommand(rootCmd *cobra.Command) {

	var passkey string
	var filename string
	var key string

	var createCmd = &cobra.Command{
		Use:   "stats",
		Short: "provides the statistic of the go keystore file",
		Long:  "displays statistics regarding the go keystore file",
		Run: func(cmd *cobra.Command, args []string) {

			ks, err := openKeyStore(filename, passkey)
			if err != nil {
				fmt.Printf("unable to open keystore(%s): %v\n", filename, err)
				return
			}

			if key == "" {
				fmt.Printf("go keystore file:\t'%s'\n", filename)
				fmt.Printf("number of items:\t %d\n", ks.Count())
				if ks.Count() > 0 {
					fmt.Printf("first key:\t\t \"%s\"\n", ks.Keys()[0])
				}

				return
			}

			if !slices.Contains(ks.Keys(), key) {
				fmt.Printf("key '%s' not found in keystore(%s)\n", key, filename)
				return
			}

			kInfo, err := ks.KeyInfo(key)
			if err != nil {
				fmt.Printf("unable to get key info: %v\n", err)
			}

			fmt.Printf("key info:\n")
			fmt.Printf("  index available:\t%v\n", kInfo.Available < 0)
			fmt.Printf("  key name:\t\t%s\n", key)
			fmt.Printf("  data length:\t\t%d bytes \t[%.2f KB]\n", binary.BigEndian.Uint32(kInfo.DataLength[:]), float64(binary.BigEndian.Uint32(kInfo.DataLength[:]))/1024.0)
			fmt.Printf("  allocated space:\t%d bytes \t[%.2f KB]\n", binary.BigEndian.Uint32(kInfo.AllocatedLength[:]), float64(binary.BigEndian.Uint32(kInfo.AllocatedLength[:]))/1024.0)
			fmt.Printf("  file position:\t%d bytes \t[%.2f KB]\n", binary.BigEndian.Uint32(kInfo.Location[:]), float64(binary.BigEndian.Uint32(kInfo.Location[:]))/1024.0)

		},
	}

	// Add flags to the config command
	createCmd.Flags().StringVar(&passkey, "pass", "", "keystore access password")
	createCmd.Flags().StringVar(&filename, "file", "", "filename for the  keystore")
	createCmd.Flags().StringVar(&key, "key", "", "get details for the key")

	_ = createCmd.MarkFlagRequired("pass")
	_ = createCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(createCmd)
}
