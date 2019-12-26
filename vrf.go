package vrf

import (
	//"fmt"

	"math/big"

	"github.com/Nik-U/pbc"
)

type VRF interface {
	//	******`*************** Main Functions *********************
	Gen(lambda int)
	Eval(x *big.Int) (*pbc.Element, []*pbc.Element)
	Verify(x *big.Int, y *pbc.Element, proof []*pbc.Element) bool
	//	********************* Main Functions *********************

	//	********************* Import/Export **********************
	SetPubKey() []*pbc.Element
	MarshalPubKey() []string
	SetSecKey() []*pbc.Element
	MarshalSecKey() []string
	SetParams() (string, *pbc.Element, int, int)
	MarshalParams() []string

	GetPubKey([]*pbc.Element)
	UnMarshalPubKey([]string) error
	GetSecKey([]*pbc.Element)
	UnMarshalSecKey([]string) error
	GetParams(*pbc.Params, *pbc.Element, int, int)
	UnMarshalParams([]string) error
	//	********************* Import/Export **********************
}

type abstractVRF struct {
	pubKey  []*pbc.Element
	secKey  []*pbc.Element
	params  *pbc.Params
	pairing *pbc.Pairing
	g       *pbc.Element
	lIn     int
	lCode   int
	typeVRF string
}

func NewVRF(typeVRF ...string) *abstractVRF {
	if len(typeVRF) > 1 {
		panic("...")
	}

	var mode string
	if len(typeVRF) == 1 {
		mode = typeVRF[0]
	} else {
		mode = "DY05"
	}
	return &abstractVRF{typeVRF: mode}
}

func (aVRF *abstractVRF) Gen(lambda uint32) {
	if aVRF.typeVRF == "" {
		panic("...")
	}
	switch aVRF.typeVRF {
	case "DY05":
		aVRF.DY05Gen(lambda)
	case "BMR10":
		aVRF.BMR10Gen(lambda)
	case "DOD03":
		aVRF.DOD03Gen(lambda)
	default:
		aVRF.DY05Gen(lambda)
	}
}

func (aVRF *abstractVRF) Eval(x *big.Int) (*pbc.Element, []*pbc.Element) {
	if aVRF.typeVRF == "" {
		panic("...")
	}
	var value *pbc.Element
	var proof []*pbc.Element

	switch aVRF.typeVRF {
	case "DY05":
		value, proof = aVRF.DY05Eval(x)
	case "BMR10":
		value, proof = aVRF.BMR10Eval(x)
	case "DOD03":
		value, proof = aVRF.DOD03Eval(x)
	default:
		value, proof = aVRF.DY05Eval(x)
	}

	return value, proof
}

func (aVRF *abstractVRF) Verify(x *big.Int, y *pbc.Element, proof []*pbc.Element) bool {
	if aVRF.typeVRF == "" {
		panic("..")
	}

	var checkBit bool
	switch aVRF.typeVRF {
	case "DY05":
		checkBit = aVRF.DY05Verify(x, y, proof)
	case "BMR10":
		checkBit = aVRF.BMR10Verify(x, y, proof)
	case "DOD03":
		checkBit = aVRF.DOD03Verify(x, y, proof)
	default:
		checkBit = aVRF.DY05Verify(x, y, proof)
	}

	return checkBit
}

func (aVRF *abstractVRF) GenNewPubKey() {
	if aVRF.typeVRF == "" {
		panic("..")
	}

	switch aVRF.typeVRF {
	case "DY05":
		aVRF.DY05GenNewPubKey()
	case "BMR10":
		aVRF.BMR10GenNewPubKey()
	case "DOD03":
		aVRF.DOD03GenNewPubKey()
	default:
		aVRF.DY05GenNewPubKey()
	}
}

func (aVRF *abstractVRF) SetPubKey(pubKey []*pbc.Element) {
	newPubKey := aVRF.MapArrayToCurve(pubKey)
	aVRF.pubKey = newPubKey
}

func (aVRF *abstractVRF) UnMarshalPubKey(pubKey []string) {
	var pubKey1 []*pbc.Element
	for i := 0; i < len(pubKey); i++ {
		element, _ := aVRF.pairing.NewG1().SetString(pubKey[i], 10)
		pubKey1 = append(pubKey1, element)
	}
	aVRF.SetPubKey(pubKey1)
}

func (aVRF *abstractVRF) SetSecKey(secKey []*pbc.Element) {
	newSecKey := aVRF.MapArrayToCurveZ(secKey)
	aVRF.secKey = newSecKey
	aVRF.GenNewPubKey()
}

func (aVRF *abstractVRF) UnMarshalSecKey(secKey []string) {
	var secKey1 []*pbc.Element
	for i := 0; i < len(secKey); i++ {
		element, _ := aVRF.pairing.NewZr().SetString(secKey[i], 10)
		secKey1 = append(secKey1, element)
	}
	aVRF.SetSecKey(secKey1)
}

func (aVRF *abstractVRF) SetParams(params string, generator []byte, lengthInput int, lengthCode int) {
	aVRF.params, _ = pbc.NewParamsFromString(params)
	aVRF.pairing = aVRF.params.NewPairing()
	aVRF.g = aVRF.pairing.NewG1().SetBytes(generator)
	aVRF.lIn = lengthInput
	aVRF.lCode = lengthCode
}

func (aVRF *abstractVRF) UnMarshalParams(allParams []string) {
	aVRF.params, _ = pbc.NewParamsFromString(allParams[0])
	aVRF.pairing = aVRF.params.NewPairing()
	aVRF.g, _ = aVRF.pairing.NewG1().SetString(allParams[1], 10)
	aVRF.lIn = 0
	aVRF.lCode = 0
}
func (aVRF *abstractVRF) GetPubKey() []*pbc.Element {
	var pubKey []*pbc.Element
	for i := 0; i < len(aVRF.pubKey); i++ {
		pubKey = append(pubKey, aVRF.pubKey[i])
	}
	return pubKey
}

func (aVRF *abstractVRF) MarshalPubKey() []string {
	var pubKey []string
	for i := 0; i < len(aVRF.pubKey); i++ {
		pubKey = append(pubKey, aVRF.pubKey[i].String())
	}
	return pubKey
}

func (aVRF *abstractVRF) GetSecKey() []*pbc.Element {
	var secKey []*pbc.Element
	for i := 0; i < len(aVRF.secKey); i++ {
		secKey = append(secKey, aVRF.secKey[i])
	}
	return secKey
}

func (aVRF *abstractVRF) MarshalSecKey() []string {
	var secKey []string
	for i := 0; i < len(aVRF.secKey); i++ {
		secKey = append(secKey, aVRF.secKey[i].String())
	}
	return secKey
}

func (aVRF *abstractVRF) GetParams() (string, []byte, int, int) {
	params := aVRF.params.String()
	generator := aVRF.g.Bytes()
	lengthInput := aVRF.lIn
	lengthCode := aVRF.lCode
	return params, generator, lengthInput, lengthCode
}

func (aVRF *abstractVRF) MarshalParams() []string {
	var allParams []string
	params := aVRF.params.String()
	allParams = append(allParams, params)
	generator := aVRF.g.String()
	allParams = append(allParams, generator)
	lengthInput := ""
	allParams = append(allParams, lengthInput)
	lengthCode := ""
	allParams = append(allParams, lengthCode)
	return allParams
}

func (aVRF *abstractVRF) MapArrayToCurve(arr []*pbc.Element) []*pbc.Element {
	var curveArr []*pbc.Element
	for i := 0; i < len(arr); i++ {
		curveArr = append(curveArr, aVRF.pairing.NewG1().SetBytes(arr[i].Bytes()))
	}
	return curveArr
}

func (aVRF *abstractVRF) MapElementToCurveT(ele *pbc.Element) *pbc.Element {
	var curveEle *pbc.Element
	curveEle = aVRF.pairing.NewGT().SetBytes(ele.Bytes())
	return curveEle
}

func (aVRF *abstractVRF) MapElementToCurve1(ele *pbc.Element) *pbc.Element {
	var curveEle *pbc.Element
	curveEle = aVRF.pairing.NewG1().SetBytes(ele.Bytes())
	return curveEle
}

func (aVRF *abstractVRF) MapArrayToCurveZ(arr []*pbc.Element) []*pbc.Element {
	var curveArr []*pbc.Element
	for i := 0; i < len(arr); i++ {
		curveArr = append(curveArr, aVRF.pairing.NewZr().SetBytes(arr[i].Bytes()))
	}
	return curveArr
}
