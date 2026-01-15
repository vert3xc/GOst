package cmd

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vert3xc/GOst/gostr1323565.1.006-2017"
	"github.com/vert3xc/GOst/streebog"
)

var rngCmd = &cobra.Command{
	Use:   "rng [seed-hex] [length]",
	Short: "GOST R 1323565.1.006-2017 pseudorandom generator",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		seed, err := hex.DecodeString(args[0])
		if err != nil {
			return fmt.Errorf("invalid seed hex")
		}

		n, err := strconv.Atoi(args[1])
		if err != nil || n <= 0 {
			return fmt.Errorf("invalid length")
		}

		h, _ := streebog.New(512)
		rng := gostr1323565_1_006_2017.New(seed, h)

		buf := make([]byte, n)
		rng.Read(buf)

		fmt.Printf("%x\n", buf)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rngCmd)
}
