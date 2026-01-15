package gostr1323565_1_006_2017

import (
	"hash"
	"math/big"

)

type GostRng struct {
	U []byte
	h int
	m int
	H hash.Hash
}

func New(K []byte, hsh hash.Hash) *GostRng {
	m := hsh.BlockSize() * 8
	h := hsh.Size() * 8
	var s int
	if h == 256{
		s = 256
	} else {
		s = 384
	}
	filling := make([]byte, (m - s - 8)/8)
	initial := append(K, filling...)
	return &GostRng{
		U: initial,
		h: h,
		m: m,
		H: hsh,
	}
}

func (r *GostRng) Read(p []byte) (int, error) {
	t := len(p) * 8
	q := t / r.h
	rem := t - q * r.h
	one := new(big.Int).SetInt64(1)
	mod := new(big.Int).Lsh(big.NewInt(1), uint(r.m-8))
	R := make([]byte, 0)
	for i := 0; i < q; i++ {
		u_num := new(big.Int).SetBytes(r.U)
		u_num.Add(u_num, one)
		u_num.Mod(u_num, mod)
		u_num.FillBytes(r.U)
		r.H.Write(r.U)
		C := r.H.Sum(nil)
		R = append(C, R...)
		r.H.Reset()
	}
	if rem != 0 {
		u_num := new(big.Int).SetBytes(r.U)
		u_num.Add(u_num, one)
		u_num.Mod(u_num, mod)
		u_num.FillBytes(r.U)
		r.H.Write(r.U)
		C := r.H.Sum(nil)
		R = append(C[len(C) - rem / 8:], R...)
	}
	copy(p, R[:len(p)])
	return len(p), nil
}