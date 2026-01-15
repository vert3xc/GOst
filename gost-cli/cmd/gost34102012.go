package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"crypto/rand"

	"github.com/spf13/cobra"
	"github.com/vert3xc/GOst/gostr34102012"
)

var gostSignCmd = &cobra.Command{
	Use:   "gost-sign [hexmessage]",
	Short: "Sign a message using GOST 34.10-2012",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyHex, _ := cmd.Flags().GetString("key")
		curveName, _ := cmd.Flags().GetString("curve")

		curve, ok := gostr34102012.Curves[curveName]
		if !ok {
			log.Fatalf("unknown curve: %s", curveName)
		}

		d, ok := new(big.Int).SetString(keyHex, 16)

		if !ok {
			log.Fatalf("hex key parsing error %s\n", keyHex)
		}

		priv := &gostr34102012.GostPrivKey{
			ParentCurve: curve,
			D:           d,
		}

		msg, err := hex.DecodeString(args[0])
		if err != nil {
			log.Fatalf("invalid message hex: %v", err)
		}

		sig, err := priv.SignMessage(rand.Reader, msg, nil)
		if err != nil {
			log.Fatalf("sign error: %v", err)
		}

		fmt.Printf("Signature: %x\n", sig)
	},
}

var gostVerifyCmd = &cobra.Command{
	Use:   "gost-verify [hexmessage] [hexsignature]",
	Short: "Verify a GOST 34.10-2012 signature",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		curveName, _ := cmd.Flags().GetString("curve")
		pubXHex, _ := cmd.Flags().GetString("pubX")
		pubYHex, _ := cmd.Flags().GetString("pubY")

		curve, ok := gostr34102012.Curves[curveName]
		if !ok {
			log.Fatalf("unknown curve: %s", curveName)
		}

		pubX, ok := new(big.Int).SetString(pubXHex, 16)
		if !ok {
			log.Fatalf("invalid pubX hex")
		}
		pubY, ok := new(big.Int).SetString(pubYHex, 16)
		if !ok {
			log.Fatalf("invalid pubY hex")
		}

		pub := &gostr34102012.GostPubKey{
			ParentCurve: curve,
			X:           pubX,
			Y:           pubY,
		}

		msg, err := hex.DecodeString(args[0])
		if err != nil {
			log.Fatalf("invalid message hex: %v", err)
		}

		sig, err := hex.DecodeString(args[1])
		if err != nil {
			log.Fatalf("invalid signature hex: %v", err)
		}

		ok = gostr34102012.VerifyMessage(pub, msg, sig)
		if ok {
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Signature is invalid")
		}
	},
}

func init() {
	gostSignCmd.Flags().String("key", "", "private key in hex")
	gostSignCmd.Flags().String("curve", "256a", "curve name (256a, 512a, 512b, 512c, test)")
	gostSignCmd.MarkFlagRequired("key")

	gostVerifyCmd.Flags().String("curve", "256a", "curve name (256a, 512a, 512b, 512c, test)")
	gostVerifyCmd.Flags().String("pubX", "", "public key X coordinate in hex")
	gostVerifyCmd.Flags().String("pubY", "", "public key Y coordinate in hex")
	gostVerifyCmd.MarkFlagRequired("pubX")
	gostVerifyCmd.MarkFlagRequired("pubY")
	rootCmd.AddCommand(gostSignCmd)
	rootCmd.AddCommand(gostVerifyCmd)
}

