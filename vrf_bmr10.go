package vrf

import (
	"math/big"
	"github.com/Nik-U/pbc"
)

func (vrf *abstractVRF) BMR10GenNewPubKey() {
	var newPubKey []*pbc.Element
	newPubKey = append(newPubKey, vrf.secKey[0])
	for i := 1; i < vrf.lCode + 1; i ++ {
		newPubKey = append(newPubKey, vrf.pairing.NewG1().PowZn(vrf.g, vrf.secKey[i]))
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
//			pk = ([r], [u]) or sk = (h, g^u[1], ..., g^u[n])
func (vrf *abstractVRF) BMR10Gen(lambda uint32) {
	// Set length
	vrf.lCode, vrf.lIn = 71, 64

	// Generate Group Parameters
	vrf.params = pbc.GenerateA(lambda, 2 * lambda)
	vrf.pairing = pbc.NewPairing(vrf.params)
	vrf.g = vrf.pairing.NewG1().Rand()

	// Generate Keys
	var pubKey, secKey []*pbc.Element
	h := vrf.pairing.NewG1().Rand()
	pubKey = append(pubKey, h)
	secKey = append(secKey, h)
	for i := 1; i < vrf.lCode + 1; i ++ {
		secKey = append(secKey, vrf.pairing.NewZr().Rand())
		pubKey = append(pubKey, vrf.pairing.NewG1().PowZn(vrf.g, secKey[i]))
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
//		v[i]: v[i-1] * g^(1/(fx[i] + u[i]))
//		c1: 1/(fx[i] + u[i])
//		c2: g^c1
// * Evaluate 3 -> value, proof
//		value: e(v[n], h)
//		proof: (v[0], v[1], ..., v[n])

func (vrf * abstractVRF) BMR10Eval(x *big.Int) (*pbc.Element, []*pbc.Element) {
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
	for i := 1; i < vrf.lIn + 1; i ++ {
		c1 := vrf.pairing.NewZr().SetInt32(int32(fx[i - 1] - '0')).ThenAdd(vrf.secKey[i]).ThenInvert()
		c2 := vrf.pairing.NewG1().PowZn(v[i - 1], c1)
		v = append(v, c2)
	}

	// Evaluate 3
	value := vrf.pairing.NewGT().Pair(v[vrf.lIn], vrf.secKey[0])
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
// * Verify1 -> check e(v[i], g^(fx[i] + u[i])) == e(v[i-1], g)
//		c1: g^fx[i] . g^(u[i])
//		c2: e(v[i], c1)
//		c3: e(v[i-1], g)
// * Verify2 -> check value = e(v[n], g)

func (vrf *abstractVRF) BMR10Verify(x *big.Int, value *pbc.Element, v []*pbc.Element) (bool) {
	value = vrf.MapElementToCurveT(value)
	v = vrf.MapArrayToCurve(v)

	// Evaluate 1
	X := PadLeft(BigToBin(x), vrf.lIn)
	fx := HCode(X)

	// Verify 1
	for i := 1; i < vrf.lIn + 1; i ++ {
		c1 := vrf.pairing.NewG1().PowZn(vrf.g, vrf.pairing.NewZr().SetInt32(int32(fx[i - 1] - '0'))).ThenMul(vrf.pubKey[i])
		c2 := vrf.pairing.NewGT().Pair(v[i], c1)
		c3 := vrf.pairing.NewGT().Pair(v[i - 1], vrf.g)
		if ! c2.Equals(c3) {
			return false
		}
	}

	// Verify 2
	if ! value.Equals(vrf.pairing.NewGT().Pair(v[vrf.lIn], vrf.pubKey[0])) {
		return false
	}
	return true
}
