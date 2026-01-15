package cmd

import "github.com/spf13/cobra"

var magmaKey string

var magmaCmd = &cobra.Command{
	Use:   "magma",
	Short: "Magma (GOST 28147-89) block cipher",
}

func init() {
	rootCmd.AddCommand(magmaCmd)

	magmaCmd.PersistentFlags().StringVarP(
		&magmaKey,
		"key",
		"k",
		"",
		"hex-encoded 256-bit key",
	)

	magmaCmd.MarkPersistentFlagRequired("key")
}