package vrf

/*
import (
	"fmt"
	"math/big"
	"github.com/Nik-U/pbc"
)

type HW10VRF struct {
	pk []*pbc.Element
	sk []*pbc.Element
	params *pbc.Params
	pairing *pbc.Pairing
	g *pbc.Element
	lec int
	lin int
}

func (vrf *HW10VRF) Gen(lambda uint32) {
	vrf.lin = 64

	vrf.params = pbc.GenerateA(lambda, 2 * lambda)
	vrf.pairing = pbc.NewPairing(vrf.params)
	vrf.g = vrf.pairing.NewG1().Rand()
	h := vrf.pairing.NewG1().Rand()
	var pk, sk []*pbc.Element
	pk = append(pk, h)
	sk = append(sk, h)
	for i := 1; i < vrf.lin + 1; i ++ {
		sk = append(sk, vrf.pairing.NewZr().Rand())
		pk = append(pk, vrf.pairing.NewG1().PowZn(vrf.g, sk[i]))
	}
	vrf.sk = sk
	vrf.pk = pk
}

func (vrf *HW10VRF) Eval(x *big.Int) (*pbc.Element, []*pbc.Element) {
	X := PadLeft(BigToBin(x), vrf.lin)
	if len(X) != vrf.lin {
		panic("...")
	}

	var v []*pbc.Element
	v = append(v, vrf.pairing.NewG1().Set(vrf.g))
	for i := 1; i < vrf.lin + 1; i ++ {
		if X[i - 1] == '1' {
			v = append(v, vrf.pairing.NewG1().PowZn(v[i - 1], vrf.sk[i]))
		} else {
			v = append(v, vrf.pairing.NewG1().Set(v[i - 1]))
		}
	}

	value := vrf.pairing.NewGT().Pair()

}*/