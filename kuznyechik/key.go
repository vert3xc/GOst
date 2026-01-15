package kuznyechik

import (
	"github.com/vert3xc/GOst/internal"
)

func F(k1, k2 *[16]byte, c [16]byte) {
    tmp := internal.Xor(*k1, c)
	S(&tmp)
    L(&tmp)
	t := *k1
	*k1 = internal.Xor(tmp, *k2)
	*k2 = t
}

func (c *kuznyechikCipher) expandKey(key []byte) {
    keyConstants := GetKeyConstants()
    var k1, k2 [16]byte
    copy(k1[:], key[:16])
    copy(k2[:], key[16:])

    c.subkeys[0] = k1
    c.subkeys[1] = k2

    for i := 0; i < 4; i++ {
        for j := 0; j < 8; j++ {
            F(&k1, &k2, keyConstants[8*i+j])
        }
        c.subkeys[2+2*i]   = k1
        c.subkeys[2+2*i+1] = k2
    }
}
