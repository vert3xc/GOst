package streebog

func add512(a, b *[64]byte) {
	var carry uint16
	for i := 63; i >= 0; i-- {
		carry = uint16(a[i]) + uint16(b[i]) + (carry >> 8)
		a[i] = byte(carry & 0xff)
	}
}

func xor512(a, b, c *[64]byte) {
	for i := 0; i < 64; i++ {
		c[i] = a[i] ^ b[i]
	}
}

func subBytes(state *[64]byte) {
  for i := 0; i < 64; i++ {
    state[i] = PI[state[i]]
  }
}

func permute(state *[64]byte) {
  var tmp [64]byte
  copy(tmp[:], state[:])
  for i := 0; i < 64; i++ {
    state[i] = tmp[TAU[i]]
  }
}

func linearTransform(state *[64]byte) {
    var input [8]uint64
    for i := 0; i < 8; i++ {
        input[i] = uint64(state[i*8])<<56 |
            uint64(state[i*8+1])<<48 |
            uint64(state[i*8+2])<<40 |
            uint64(state[i*8+3])<<32 |
            uint64(state[i*8+4])<<24 |
            uint64(state[i*8+5])<<16 |
            uint64(state[i*8+6])<<8 |
            uint64(state[i*8+7])
    }

    var output [8]uint64
    for i := 0; i < 8; i++ {
        for j := 0; j < 64; j++ {
            if (input[i]>>j)&1 == 1 {
                output[i] ^= A[63-j]
            }
        }
    }
    for i := 0; i < 8; i++ {
        state[i*8] = byte(output[i] >> 56)
        state[i*8+1] = byte(output[i] >> 48)
        state[i*8+2] = byte(output[i] >> 40)
        state[i*8+3] = byte(output[i] >> 32)
        state[i*8+4] = byte(output[i] >> 24)
        state[i*8+5] = byte(output[i] >> 16)
        state[i*8+6] = byte(output[i] >> 8)
        state[i*8+7] = byte(output[i])
    }
}

func keySchedule(keys *[64]byte, iter int) {
	xor512(keys, &C[iter], keys)
	subBytes(keys)
	permute(keys)
	linearTransform(keys)
}

func E(keys, block, state *[64]byte) {
	xor512(block, keys, state)
	for i := 0; i < 12; i++ {
		subBytes(state)
		permute(state)
		linearTransform(state)
		keySchedule(keys, i)
		xor512(state, keys, state)
	}
}

func G(n, hash, m *[64]byte) {
	var keys [64]byte
	var tmp [64]byte

	xor512(n, m, &keys)

	subBytes(&keys)
	permute(&keys)
	linearTransform(&keys)

	E(&keys, hash, &tmp)

	xor512(&tmp, n, &tmp)
	xor512(&tmp, hash, n)
}