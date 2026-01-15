package internal

func Xor(a, b [16]byte) [16]byte{
	out := [16]byte{}
	for i := 0; i < 16; i++ {
		out[i] = a[i] ^ b[i]
	}
	return out
}

func MulGF(a, b byte) byte {
    var res byte
    for i := 0; i < 8; i++ {
        if (b & 1) != 0 {
            res ^= a
        }
        hiBitSet := (a & 0x80) != 0
        a <<= 1
        if hiBitSet {
            a ^= 0xC3
        }
        b >>= 1
    }
    return res
}