package gostr34102012

import (
	"crypto"
	"math/big"
	"io"
	"slices"
	"fmt"

	"github.com/vert3xc/GOst/streebog"
)

type GostPrivKey struct {
	ParentCurve *Curve
	D *big.Int
}

type GostPubKey struct {
	ParentCurve *Curve
	X *big.Int
	Y *big.Int
}

func (pr *GostPrivKey) Public() *GostPubKey {
	pubx, puby := pr.ParentCurve.ScalarMult(pr.ParentCurve.Gx, pr.ParentCurve.Gy, pr.D)
	return &GostPubKey{
		ParentCurve: pr.ParentCurve,
		X: pubx,
		Y: puby,
	}
}


func (pr *GostPrivKey) Equal(x *GostPrivKey) bool {
    if pr == x {
        return true
    }
    if pr == nil || x == nil {
        return false
    }
    return pr.ParentCurve == x.ParentCurve &&
        pr.D.Cmp(x.D) == 0
}

func (pub *GostPubKey) Equal(x *GostPubKey) bool {
    if pub == x {
        return true
    }
    if pub == nil || x == nil {
        return false
    }
    return pub.ParentCurve == x.ParentCurve &&
        pub.X.Cmp(x.X) == 0 &&
        pub.Y.Cmp(x.Y) == 0
}

func (pr *GostPrivKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	rev := slices.Clone(digest)
	slices.Reverse(rev)
	e := new(big.Int).SetBytes(rev)
	e.Mod(e, pr.ParentCurve.Q)
	if e.Sign() == 0 {
		e.SetInt64(1)
	}
	r := new(big.Int).SetInt64(0)
	s := new(big.Int).SetInt64(0)
	for ; r.Sign() == 0 || s.Sign() == 0; {
		nonce := make([]byte, pr.ParentCurve.BitSize / 8)
		_, err := rand.Read(nonce)
		if err != nil {
			return nil, err
		}
		k := new(big.Int).SetBytes(nonce)
		k.Mod(k, pr.ParentCurve.Q)
		r, _ = pr.ParentCurve.ScalarMult(pr.ParentCurve.Gx, pr.ParentCurve.Gy, k)
		r.Mod(r, pr.ParentCurve.Q)
		s.Add(new(big.Int).Mul(r, pr.D), new(big.Int).Mul(k, e)).Mod(s, pr.ParentCurve.Q)
	}
	size := (pr.ParentCurve.BitSize + 7) / 8
    sig := make([]byte, 2*size)

    s.FillBytes(sig[:size])
    r.FillBytes(sig[size:])
	return sig, nil
}

func (pr *GostPrivKey) SignMessage(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	hsh, _ := streebog.New(512)
	hsh.Write(digest)
	hashdigest := hsh.Sum(nil)
	return pr.Sign(rand, hashdigest, opts)
}

func (pr *GostPrivKey) TestSign(e, k *big.Int) (*big.Int, *big.Int, error) {
	
	r := new(big.Int).SetInt64(0)
	s := new(big.Int).SetInt64(0)


	r, _ = pr.ParentCurve.ScalarMult(pr.ParentCurve.Gx, pr.ParentCurve.Gy, k)
	r.Mod(r, pr.ParentCurve.Q)
	s.Add(new(big.Int).Mul(r, pr.D), new(big.Int).Mul(k, e)).Mod(s, pr.ParentCurve.Q)


	return s, r, nil
}

func Verify(pub *GostPubKey, digest, sig []byte) bool {
	size := (pub.ParentCurve.BitSize + 7) / 8
	r := new(big.Int).SetBytes(sig[size:])
	s := new(big.Int).SetBytes(sig[:size])
	e := new(big.Int).SetBytes(digest)
	e.Mod(e, pub.ParentCurve.Q)
	if e.Sign() == 0 {
		e.SetInt64(1)
	}
	if r.Sign() <= 0 || r.Cmp(pub.ParentCurve.Q) >= 0 || s.Sign() <= 0 || s.Cmp(pub.ParentCurve.Q) >= 0 {
		fmt.Printf("r or s not in bound\n")
		fmt.Printf("R:%d\nS:%d\n", r, s)
    	return false
	}
	v := new(big.Int).ModInverse(e, pub.ParentCurve.Q)
	z1 := new(big.Int).Mul(s, v)
	z1.Mod(z1, pub.ParentCurve.Q)
	z2 := new(big.Int).Mul(r, v)
	z2.Neg(z2).Mod(z2, pub.ParentCurve.Q)
	x0, y0 := pub.ParentCurve.ScalarMult(pub.ParentCurve.Gx, pub.ParentCurve.Gy, z1)
	x1, y1 := pub.ParentCurve.ScalarMult(pub.X, pub.Y, z2)
	R, _ := pub.ParentCurve.Add(x0, y0, x1, y1)
	R.Mod(R, pub.ParentCurve.Q)
	return R.Cmp(r) == 0
}

func VerifyMessage(pub *GostPubKey, digest, sig []byte) bool {
	hsh, _ := streebog.New(512)
	hsh.Write(digest)
	hashdigest := hsh.Sum(nil)
	rev := slices.Clone(hashdigest)
	slices.Reverse(rev)
	return Verify(pub, rev, sig)
}

func TestVerify(pub *GostPubKey, e, r, s *big.Int) bool {
	if e.Sign() == 0 {
		e.SetInt64(1)
	}
	if r.Sign() <= 0 || r.Cmp(pub.ParentCurve.Q) >= 0 || s.Sign() <= 0 || s.Cmp(pub.ParentCurve.Q) >= 0 {
    	return false
	}
	v := new(big.Int).ModInverse(e, pub.ParentCurve.Q)
	z1 := new(big.Int).Mul(s, v)
	z1.Mod(z1, pub.ParentCurve.Q)
	z2 := new(big.Int).Mul(r, v)
	z2.Neg(z2).Mod(z2, pub.ParentCurve.Q)
	x0, y0 := pub.ParentCurve.ScalarMult(pub.ParentCurve.Gx, pub.ParentCurve.Gy, z1)
	x1, y1 := pub.ParentCurve.ScalarMult(pub.X, pub.Y, z2)
	R, _ := pub.ParentCurve.Add(x0, y0, x1, y1)
	return R.Cmp(r) == 0
}