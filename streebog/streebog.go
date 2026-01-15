package streebog

import (
	"hash"
	"errors"
)

type streebogHash struct {
	state         [64]byte
	n             [64]byte
	sigma         [64]byte
	x             [64]byte
	nx 	          int
	blocksize     int
}

func New(blocksize int) (hash.Hash, error) {
	if blocksize != 256 && blocksize != 512 {
		return nil, errors.New("streebog: invalid blocksize (must be 256 or 512 bits)")
	}
	s := &streebogHash{blocksize: blocksize}
	s.Reset()
	return s, nil
}

func (s *streebogHash) Write(p []byte) (nu int, err error) {
	nu = len(p)
	if s.nx > 0{
		n := copy(s.x[s.nx:], p)
		s.nx += n
		if s.nx == 64 {
			block(s, s.x[:])
			s.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= 64 {
		n := len(p) &^ (64 - 1)
		block(s, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		s.nx = copy(s.x[:], p)
	}
	return
}

func (s *streebogHash) checkSum() []byte {
	if s.nx > 0 {
		block(s, s.x[:s.nx])
	}
	finalize(s)
	return s.state[:s.blocksize/8]
}

func (s *streebogHash) Sum(b []byte) []byte {
	tmp := *s
	hash := tmp.checkSum()
	return append(b, hash...)
}

func (s *streebogHash) Size() int{
	if s.blocksize == 256{
		return 32
	} else {
		return 64
	}
}

func (s *streebogHash) BlockSize() int{
	return 64
}

func (s *streebogHash) Reset() {
    if s.blocksize == 256 {
        for i := range s.state {
            s.state[i] = 1
			s.n[i] = 0
			s.sigma[i] = 0
        }
    } else {
        for i := range s.state {
            s.state[i] = 0
			s.n[i] = 0
			s.sigma[i] = 0
        }
    }
}