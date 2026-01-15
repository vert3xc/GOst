package kuznyechik

import "github.com/vert3xc/GOst/internal"

var lVec = [16]byte{1, 148, 32, 133, 16, 194, 192, 1, 251, 1, 192, 194, 16, 133, 32, 148}

func shiftRight(block *[16]byte) {
    tmp := block[15]
    for i := 15; i > 0; i-- {
        block[i] = block[i-1]
    }
    block[0] = tmp
}

func shiftLeft(block *[16]byte) {
    tmp := block[0]
    for i := 0; i < 15; i++ {
        block[i] = block[i+1]
    }
    block[15] = tmp
}

func R(a *[16]byte) {
	var out byte
	shiftRight(a)
	for i, b := range lVec{
		out ^= internal.MulGF(a[i], b)
	}
	a[0] = out
}

func InvR(a *[16]byte) {
    var out byte
	
	for i, b := range lVec{
		out ^= internal.MulGF(a[i], b)
	}
	a[0] = out
    shiftLeft(a)
}

func L(a *[16]byte) {
	for i := 0; i < 16; i++ {
		R(a)
	}
}

func InvL(a *[16]byte) {
    for i := 0; i < 16; i++ {
        InvR(a)
    }
}

func S(a *[16]byte) {
	for i := 0; i < 16; i++ {
		a[i] = KuznyechikSBox[a[i]]
	}
}

func InvS(a *[16]byte) {
	for i := 0; i < 16; i++ {
		a[i] = KuznyechikInvSBox[a[i]]
	}
}

func GetKeyConstants() [32][16]byte {
    var constants [32][16]byte
    for i := 1; i <= 32; i++ {
        var curConst [16]byte
        curConst[15] = byte(i)
        L(&curConst)
        constants[i-1] = curConst
    }
    return constants
}