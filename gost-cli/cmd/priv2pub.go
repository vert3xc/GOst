package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/vert3xc/GOst/gostr34102012"
)

var priv2pubCmd = &cobra.Command{
	Use:   "priv2pub [hexprivate]",
	Short: "Derive public key from private key (GOST 34.10-2012)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyHex := args[0]
		curveName, _ := cmd.Flags().GetString("curve")

		curve, ok := gostr34102012.Curves[curveName]
		if !ok {
			log.Fatalf("unknown curve: %s", curveName)
		}

		d, ok := new(big.Int).SetString(keyHex, 16)
		if !ok {
			log.Fatalf("invalid private key hex")
		}

		priv := &gostr34102012.GostPrivKey{
			ParentCurve: curve,
			D:           d,
		}

		pub := priv.Public()
		fmt.Printf("Public key coordinates:\nX: %x\nY: %x\n", pub.X, pub.Y)
	},
}

func init() {
	priv2pubCmd.Flags().String("curve", "256a", "curve name (256a, 512a, 512b, 512c, test)")
	rootCmd.AddCommand(priv2pubCmd)
}