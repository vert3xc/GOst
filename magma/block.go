package magma

import (
    "github.com/vert3xc/GOst/internal"
    "math/bits"
)

func t(a uint32) uint32 {
    var out uint32
    for i := 0; i < 8; i++ {
        nibble := byte((a >> (4 * i)) & 0xF)
        sub := MagmaSBox[i][nibble]
        out |= uint32(sub) << (4 * i)
    }
    return out
}

func g(a, x uint32) uint32 {
    v := a + x
    v = t(v)
    return bits.RotateLeft32(v, 11)
}

func (c *magmaCipher) EncryptBlock(src []byte) []byte {
    dst := make([]byte, 8)
    A := internal.BEToUint32(src[:4])
    B := internal.BEToUint32(src[4:])
    for i := 0; i < 32; i++ {
        t := B
        B = A ^ g(B, c.subkeys[i])
        A = t
    }

    internal.Uint32ToBE(dst[4:], A)
    internal.Uint32ToBE(dst[:4], B)
    return dst
}

func (c *magmaCipher) DecryptBlock(src []byte) []byte {
    dst := make([]byte, 8)
    A := internal.BEToUint32(src[:4])
    B := internal.BEToUint32(src[4:])
    for i := 0; i < 32; i++ {
        t := B
        B = A ^ g(B, c.subkeys[31-i])
        A = t
    }

    internal.Uint32ToBE(dst[4:], A)
    internal.Uint32ToBE(dst[:4], B)
    return dst
}

