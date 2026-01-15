package cmd

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/vert3xc/GOst/kuznyechik"
)

var kuznyechikDecryptCmd = &cobra.Command{
	Use:   "decrypt [hexdata]",
	Short: "Decrypt data using Kuznyechik",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := decodeHexArg(kuznyechikKey)
		if err != nil {
			return err
		}

		block, err := kuznyechik.NewCipher(key)
		if err != nil {
			return err
		}

		var data []byte
		if len(args) == 1 {
			data, err = decodeHexArg(args[0])
		} else {
			data, err = io.ReadAll(os.Stdin)
		}
		if err != nil {
			return err
		}

		if len(data)%block.BlockSize() != 0 {
			return fmt.Errorf(
				"input length (%d bytes) is not a multiple of block size (%d bytes)",
				len(data),
				block.BlockSize(),
			)
		}

		out := make([]byte, len(data))
		for i := 0; i < len(data); i += block.BlockSize() {
			block.Decrypt(out[i:], data[i:])
		}

		cmd.Println(hex.EncodeToString(out))
		return nil
	},
}

func init() {
	kuznyechikCmd.AddCommand(kuznyechikDecryptCmd)
}
