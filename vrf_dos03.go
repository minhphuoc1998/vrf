package vrf

import (
	"math/big"

	"github.com/Nik-U/pbc"
)

func (vrf *abstractVRF) DOD03GenNewPubKey() {
	var newPubKey []*pbc.Element
	newPubKey = append(newPubKey, vrf.secKey[0])
	for i := 1; i < vrf.lCode+1; i++ {
		newPubKey = append(newPubKey, vrf.pairing.NewG1().PowZn(vrf.secKey[0], vrf.secKey[i]))
	}
	vrf.pubKey = newPubKey
}

// ****** Generation ******
// - In: lambda
// - Out: None
// * Set length
//		lIn: length of input
//		lCode: length of code
// * Generate Group Parameters
// 		params: group parameters
// 		pairing: pair in group
// 		g: group generator
//	* Generate Keys
// 		secKey: secret key
// 			sk = ([r], u) or sk = (h, u[1], ..., u[n]) where n = lCode
// 		pubKey: public key
//			pk = ([r], [u]) or sk = (h, h^u[1], ..., h^u[n])
func (vrf *abstractVRF) DOD03Gen(lambda uint32) {
	// Set length
	vrf.lCode, vrf.lIn = 71, 64

	// Generate Group Parameters
	vrf.params = pbc.GenerateA(lambda, 2*lambda)
	vrf.pairing = pbc.NewPairing(vrf.params)
	vrf.g = vrf.pairing.NewG1().Rand()

	// Generate Keys
	h := vrf.pairing.NewG1().Rand()
	var secKey, pubKey []*pbc.Element
	secKey = append(secKey, h)
	pubKey = append(pubKey, h)
	for i := 1; i < vrf.lCode+1; i++ {
		secKey = append(secKey, vrf.pairing.NewZr().Rand())
		pubKey = append(pubKey, vrf.pairing.NewG1().PowZn(h, secKey[i]))
	}
	vrf.secKey = secKey
	vrf.pubKey = pubKey
}

// ***** Evaluation ******
// - In:
//		x: seed
// - Out:
//		value: value
//		proof: proof
// * Evaluate 1 -> encode x
//		X: binary of x
//		fx: code(X)
// * Evaluate 2 -> value, proof
//		v[i]: v[i-1] * u[i] if fx[i] == 1 else v[i-1]
// * Evaluate 3 -> value, proof
//		value: v[n]
//		proof: (v[0], v[1], ..., v[n])
func (vrf *abstractVRF) DOD03Eval(x *big.Int) (*pbc.Element, []*pbc.Element) {
	// Evaluate 1
	X := PadLeft(BigToBin(x), vrf.lIn)
	if len(X) != vrf.lIn {
		panic("...")
	}
	fx := HCode(X)
	if len(fx) != vrf.lCode {
		panic("...")
	}
	// Evaluate 2
	var v []*pbc.Element
	v = append(v, vrf.pairing.NewG1().Set(vrf.g))
	for i := 1; i < vrf.lCode+1; i++ {
		if fx[i-1] == '1' {
			v = append(v, vrf.pairing.NewG1().PowZn(v[i-1], vrf.secKey[i]))
		} else {
			v = append(v, vrf.pairing.NewG1().Set(v[i-1]))
		}
	}
	// Evaluate 3
	value := v[vrf.lCode]
	proof := v
	return value, proof
}

// ***** Verification *****
// - In:
//		x: seed
//		value: value
//		proof: proof
// - Out:
//		0/1 or valid/invalid
// * Evaluate1 -> encode x
//		X: binary of x
//		fx: code(X)
// * Verify -> check e(v[i], h) == e(v[i-1], h^u[i] if fx[i] == 1 else h)
//		c1: e(v[i-1], h^u[i] if fx[i] == 1 else h)
//		c2: e(v[i], h)
func (vrf *abstractVRF) DOD03Verify(x *big.Int, y *pbc.Element, v []*pbc.Element) bool {
	y = vrf.MapElementToCurveT(y)
	v = vrf.MapArrayToCurve(v)

	// Evaluate 1
	X := PadLeft(BigToBin(x), vrf.lIn)
	if len(X) != vrf.lIn {
		panic("...")
	}
	fx := HCode(X)
	if len(fx) != vrf.lCode {
		panic("...")
	}

	// Verify
	for i := 1; i < vrf.lCode+1; i++ {
		var c1 *pbc.Element
		if fx[i-1] == '1' {
			c1 = vrf.pairing.NewGT().Pair(v[i-1], vrf.pubKey[i])
		} else {
			c1 = vrf.pairing.NewGT().Pair(v[i-1], vrf.pubKey[0])
		}
		c2 := vrf.pairing.NewGT().Pair(v[i], vrf.pubKey[0])
		if !c1.Equals(c2) {
			return false
		}
	}
	return true
}
