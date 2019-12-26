package vrf

import (
	"math/big"

	"github.com/Nik-U/pbc"
)

func (vrf *abstractVRF) DY05GenNewPubKey() {
	var newPubKey []*pbc.Element
	newPubKey = append(newPubKey, vrf.pairing.NewG1().PowZn(vrf.g, vrf.secKey[0]))
	vrf.pubKey = newPubKey
}

// ****** Generation ******
// - In: lambda
// - Out: None
// Generate Group Parameters
// 		params: group parameters
// 		pairing: pair in group
// 		g: group generator
//	Generate Keys
// 		secKey: secret key
// 			sk = r
// 		pubKey: public key
//			pk = [r] or sk = g^r
func (vrf *abstractVRF) DY05Gen(lambda uint32) {
	// Generate Group Parameters
	vrf.params = pbc.GenerateA(lambda, 2*lambda)
	vrf.pairing = pbc.NewPairing(vrf.params)
	vrf.g = vrf.pairing.NewG1().Rand()

	// Generate Keys
	var secKey, pubKey []*pbc.Element
	secKey = append(secKey, vrf.pairing.NewZr().Rand())
	vrf.secKey = secKey
	pubKey = append(pubKey, vrf.pairing.NewG1().PowZn(vrf.g, vrf.secKey[0]))
	vrf.pubKey = pubKey
}

// ***** Evaluation ******
// - In:
//		x: seed
// - Out:
//		value: value
//		proof: proof
// Evaluate 1 -> [1/(X+r)] or g^{1/(X+r)}
//		X: x to Zr
//		t: 1/(X+r)
//		gt: g^t
// Evaluate 2 -> value, proof
//		value: e(g, gt)
//		proof: gt

func (vrf *abstractVRF) DY05Eval(x *big.Int) (*pbc.Element, []*pbc.Element) {
	// Evaluate 1
	X := vrf.pairing.NewZr().SetBig(x)
	t := vrf.pairing.NewZr().Add(X, vrf.secKey[0]).ThenInvert()
	gt := vrf.pairing.NewG1().PowZn(vrf.g, t)

	// Evaluate 2
	var value *pbc.Element
	var proof []*pbc.Element
	value = vrf.pairing.NewGT().Pair(vrf.g, gt)
	proof = append(proof, gt)

	return value, proof
}

// ***** Verification *****
// - In:
//		x: seed
//		value: value
//		proof: proof
// - Out:
//		0/1 or valid/invalid
// Verify1 -> check e(g^x . g^r, g^(1/(x+r))) == e(g, g)
//		X: x to Zr
//		gx: g^x
//		c1: e(g^x . g^r, g^(1/(x+r)))
//		c2: e(g, g)
// Verify2 -> check value == e(g^(1/(x+r)), g)
//		gt: g^(9)1/(x+r))
//		c3: e(e(g^(1/(x+r)), g))

func (vrf *abstractVRF) DY05Verify(x *big.Int, value *pbc.Element, proof []*pbc.Element) bool {
	value = vrf.MapElementToCurveT(value)
	proof = vrf.MapArrayToCurve(proof)

	X := vrf.pairing.NewZr().SetBig(x)
	// Verify 1
	gx := vrf.pairing.NewG1().PowZn(vrf.g, X)
	c1 := vrf.pairing.NewGT().Pair(vrf.pairing.NewG1().Mul(gx, vrf.pubKey[0]), proof[0])
	c2 := vrf.pairing.NewGT().Pair(vrf.g, vrf.g)
	if !(c1.Equals(c2)) {
		return false
	}

	// Verify 2
	gt := proof[0]
	c3 := vrf.pairing.NewGT().Pair(gt, vrf.g)
	if !(c3.Equals(value)) {
		return false
	}
	return true
}
