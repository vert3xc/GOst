package cmd

import "github.com/spf13/cobra"

var kuznyechikKey string

var kuznyechikCmd = &cobra.Command{
	Use:   "kuznyechik",
	Short: "Kuznyechik (GOST R 34.12-2015) block cipher",
}

func init() {
	rootCmd.AddCommand(kuznyechikCmd)

	kuznyechikCmd.PersistentFlags().StringVarP(
		&kuznyechikKey,
		"key",
		"k",
		"",
		"hex-encoded 256-bit key",
	)

	kuznyechikCmd.MarkPersistentFlagRequired("key")
}
