package kuznyechik

import (
	"crypto/cipher"
	"errors"
)

type kuznyechikCipher struct {
    subkeys [10][16]byte
	blocksize int
}

func NewCipher(key []byte) (cipher.Block, error) {
    if len(key) != 32 {
        return nil, errors.New("kuznyechik: invalid key size (must be 256 bits)")
    }
    c := &kuznyechikCipher{}
	c.blocksize = 16
    c.expandKey(key)
    return c, nil
}

func (c *kuznyechikCipher) BlockSize() int { return c.blocksize }

func (c *kuznyechikCipher) Encrypt(dst, src []byte) {
    if len(src) < 16 || len(dst) < 16 {
        panic("kuznyechik: input not full block")
    }
    res := c.EncryptBlock(src)
    copy(dst, res[:])
}

func (c *kuznyechikCipher) Decrypt(dst, src []byte) {
    if len(src) < 16 || len(dst) < 16 {
        panic("kuznyechik: input not full block")
    }
    res := c.DecryptBlock(src)
    copy(dst, res[:])
}