package kuznyechik

import (
	"github.com/vert3xc/GOst/internal"
)

func (c *kuznyechikCipher) EncryptBlock(src []byte) [16]byte {
	var tmp [16]byte
	copy(tmp[:], src)
    for i := 0; i < 9; i++ {
        tmp = internal.Xor(tmp, c.subkeys[i])
		S(&tmp)
		L(&tmp)
    }
	tmp = internal.Xor(tmp, c.subkeys[9])

    return tmp
}

func (c *kuznyechikCipher) DecryptBlock(src []byte) [16]byte {
	var tmp [16]byte
	copy(tmp[:], src)
	tmp = internal.Xor(tmp, c.subkeys[9])
    for i := 0; i < 9; i++ {
		InvL(&tmp)
		InvS(&tmp)
		tmp = internal.Xor(tmp, c.subkeys[8 - i])
    }
    return tmp
}
