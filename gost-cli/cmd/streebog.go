package cmd

import (
	"fmt"
	"io"
	"os"
	"encoding/hex"
	"log"

	"github.com/spf13/cobra"
	"github.com/vert3xc/GOst/streebog"
)

var streebogBits int

var streebogCmd = &cobra.Command{
	Use:   "streebog [hexstring]",
	Short: "Compute Streebog (GOST R 34.11-2012) hash",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := streebog.New(streebogBits)
		if err != nil {
			return err
		}
		if len(args) == 1 {
			data_bytes, err := hex.DecodeString(args[0])
			if err != nil {
				log.Fatalf("invalid input hex: %v", err)
			}
			_, err = h.Write(data_bytes)
		} else {
			_, err = io.Copy(h, os.Stdin)
		}
		if err != nil {
			return err
		}

		sum := h.Sum(nil)

		fmt.Printf("%x\n", sum)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(streebogCmd)

	streebogCmd.Flags().IntVarP(
		&streebogBits,
		"bits",
		"b",
		512,
		"hash size: 256 or 512",
	)
}
