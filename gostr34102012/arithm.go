package gostr34102012

import (
	"math/big"
)

func isInfinity(x, y *big.Int) bool {
    return x == nil && y == nil
}

func (c *Curve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
    if isInfinity(x1, y1) {
        if isInfinity(x2, y2) {
            return nil, nil
        }
        return new(big.Int).Set(x2), new(big.Int).Set(y2)
    }

    if isInfinity(x2, y2) {
        return new(big.Int).Set(x1), new(big.Int).Set(y1)
    }

    p := c.P

    if x1.Cmp(x2) == 0 {
        ySum := new(big.Int).Add(y1, y2)
        ySum.Mod(ySum, p)
        if ySum.Sign() == 0 {
            return nil, nil
        }
    }

    var lambda *big.Int

    if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 {
        num := new(big.Int).Mul(x1, x1)
        num.Mul(num, big.NewInt(3))
        num.Add(num, c.A)
        num.Mod(num, p)

        den := new(big.Int).Mul(y1, big.NewInt(2))
        den.Mod(den, p)

        denInv := new(big.Int).ModInverse(den, p)
        if denInv == nil {
            return nil, nil
        }

        lambda = num.Mul(num, denInv)
        lambda.Mod(lambda, p)

    } else {
        num := new(big.Int).Sub(y2, y1)
        num.Mod(num, p)

        den := new(big.Int).Sub(x2, x1)
        den.Mod(den, p)

        denInv := new(big.Int).ModInverse(den, p)
        if denInv == nil {
            return nil, nil
        }

        lambda = num.Mul(num, denInv)
        lambda.Mod(lambda, p)
    }

    x3 := new(big.Int).Mul(lambda, lambda)
    x3.Sub(x3, x1)
    x3.Sub(x3, x2)
    x3.Mod(x3, p)

    y3 := new(big.Int).Sub(x1, x3)
    y3.Mul(lambda, y3)
    y3.Sub(y3, y1)
    y3.Mod(y3, p)

    return x3, y3
}

func (c *Curve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
    p := c.P
    num := new(big.Int).Mul(x1, x1)
    num.Mul(num, big.NewInt(3))
    num.Add(num, c.A)
    num.Mod(num, p)

    den := new(big.Int).Mul(y1, big.NewInt(2))
    den.Mod(den, p)

    denInv := new(big.Int).ModInverse(den, p)
    if denInv == nil {
        return nil, nil
    }

    lambda := num.Mul(num, denInv)
    lambda.Mod(lambda, p)

    x3 := new(big.Int).Mul(lambda, lambda)
    x3.Sub(x3, new(big.Int).Mul(big.NewInt(2), x1))
    x3.Mod(x3, p)

    y3 := new(big.Int).Sub(x1, x3)
    y3.Mul(lambda, y3)
    y3.Sub(y3, y1)
    y3.Mod(y3, p)

    return x3, y3
}

func (c *Curve) ScalarMult(x1, y1, k *big.Int) (*big.Int, *big.Int) {
	r0_x := new(big.Int).Set(x1)
	r0_y := new(big.Int).Set(y1)
	r1_x, r1_y := c.Double(x1, y1)
	n := k.BitLen()
	for i := n - 2; i > -1; i-- {
		r0_x_copy := new(big.Int).Set(r0_x)
		r0_y_copy := new(big.Int).Set(r0_y)
		r1_x_copy := new(big.Int).Set(r1_x)
		r1_y_copy := new(big.Int).Set(r1_y)
		if k.Bit(i) == 0 {
			r0_x, r0_y = c.Double(r0_x_copy, r0_y_copy)
			r1_x, r1_y = c.Add(r0_x_copy, r0_y_copy, r1_x_copy, r1_y_copy)
		} else {
			r0_x, r0_y = c.Add(r0_x_copy, r0_y_copy, r1_x_copy, r1_y_copy)
			r1_x, r1_y = c.Double(r1_x_copy, r1_y_copy)
		}
	}
	return r0_x, r0_y
}