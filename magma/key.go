package magma

import (
    "github.com/vert3xc/GOst/internal"
)

func (c *magmaCipher) expandKey(key []byte) {
    var keywords [8]uint32
    for i := 0; i < 8; i++ {
        keywords[i] = internal.BEToUint32(key[i * 4 : (i + 1) * 4])
    }
    for i := 0; i < 24; i++ {
        c.subkeys[i] = keywords[i % 8]
    }
	for i := 24; i < 32; i++ {
		c.subkeys[i] = keywords[7 - (i - 24)]
	}
}
