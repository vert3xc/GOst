package magma

import (
	"crypto/cipher"
	"errors"
)

type magmaCipher struct {
    subkeys [32]uint32
	blocksize int
}

func NewCipher(key []byte) (cipher.Block, error) {
    if len(key) != 32 {
        return nil, errors.New("magma: invalid key size (must be 256 bits)")
    }
    c := &magmaCipher{}
	c.blocksize = 8
    c.expandKey(key)
    return c, nil
}

func (c *magmaCipher) BlockSize() int { return c.blocksize }

func (c *magmaCipher) Encrypt(dst, src []byte) {
    if len(src) < 8 || len(dst) < 8 {
        panic("magma: input not full block")
    }
    res := c.EncryptBlock(src)
    copy(dst, res)
}

func (c *magmaCipher) Decrypt(dst, src []byte) {
    if len(src) < 8 || len(dst) < 8 {
        panic("magma: input not full block")
    }
    res := c.DecryptBlock(src)
    copy(dst, res)
}