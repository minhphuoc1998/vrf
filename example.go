package vrf

import (
	"fmt"
	"math/big"
)

func Example() {
	vrf := NewVRF("BMR10")
	vrf.Gen(128)
	value, proof := vrf.Eval(big.NewInt(100))
	fmt.Println(vrf.Verify(big.NewInt(100), value, proof))

	vrf = NewVRF("DOD03")
	vrf.Gen(128)
	value, proof = vrf.Eval(big.NewInt(100))
	fmt.Println(vrf.Verify(big.NewInt(100), value, proof))

	vrf = NewVRF("DY05")
	vrf.Gen(128)
	value, proof = vrf.Eval(big.NewInt(100))
	fmt.Println(vrf.Verify(big.NewInt(100), value, proof))
}

func SampleProtocol() {
	seed := big.NewInt(123)
	// Alice
	AliceVRF := NewVRF("DY05")
	AliceVRF.Gen(128)
	params, generator, lIn, lCode := AliceVRF.GetParams()
	AlicePubKey := AliceVRF.GetPubKey()
	value, proof := AliceVRF.Eval(seed)
	// Bob
	BobVRF := NewVRF("DY05")
	BobVRF.SetParams(params, generator, lIn, lCode)
	BobVRF.SetPubKey(AlicePubKey)
	checkBit := BobVRF.Verify(seed, value, proof)
	fmt.Println(checkBit)
}

func SampleGame() {
	fmt.Println("--------------Step 1: Generation Process--------------")
	player1 := NewVRF("DY05")
	player1.Gen(128)
	p1 := player1.MarshalPubKey()
	s1 := player1.MarshalSecKey()
	params1 := player1.MarshalParams()
	fmt.Println("Player 1:")
	fmt.Println("Public Key:")
	fmt.Println("p1:", p1)
	fmt.Println("Secret Key:")
	fmt.Println("s1:", s1)

	player2 := NewVRF("DY05")
	player2.Gen(128)
	p2 := player2.MarshalPubKey()
	s2 := player2.MarshalSecKey()
	params2 := player2.MarshalParams()
	fmt.Println("Player 2:")
	fmt.Println("Public Key:")
	fmt.Println("p2:", p2)
	fmt.Println("Secret Key:")
	fmt.Println("s2:", s2)

	banker := NewVRF("DY05")
	banker.Gen(128)
	pb := banker.MarshalPubKey()
	sb := banker.MarshalSecKey()
	paramsb := banker.MarshalParams()
	fmt.Println("Banker:")
	fmt.Println("Public Key:")
	fmt.Println("pb:", pb)
	fmt.Println("Secret Key:")
	fmt.Println("sb:", sb)

	fmt.Println("--------------Step 2: Publish Key--------------")
	bankerVRF1 := NewVRF("DY05")
	bankerVRF1.UnMarshalParams(params1)
	bankerVRF1.UnMarshalPubKey(p1)
	bankerVRF2 := NewVRF("DY05")
	bankerVRF2.UnMarshalParams(params2)
	bankerVRF2.UnMarshalPubKey(p2)
	playerVRF := NewVRF("DY05")
	playerVRF.UnMarshalParams(paramsb)
	playerVRF.UnMarshalPubKey(pb)

	fmt.Println("--------------Step 3: Betting--------------")
	v1 := big.NewInt(100)
	fmt.Println("Player 1 value:", v1)
	v2 := big.NewInt(200)
	fmt.Println("Player 2 value:", v2)

	fmt.Println("--------------Step 4: Seed Generation--------------")
	seed1 := big.NewInt(123456)
	seed2, proof2 := banker.Eval(seed1)
	fmt.Println("Seed:")
	fmt.Println("seed2:", seed2)
	fmt.Println("proof2:", proof2)
	vers := playerVRF.Verify(seed1, seed2, proof2)
	fmt.Println("Verification Result:", vers)

	fmt.Println("--------------Step 5: Evaluation--------------")
	V1, P1 := player1.Eval(seed2.X())
	fmt.Println("Player 1 evaluation:")
	fmt.Println("V1:", V1)
	fmt.Println("P1:", P1)
	V2, P2 := player2.Eval(seed2.X())
	fmt.Println("Player 2 evaluation:")
	fmt.Println("V2:", V2)
	fmt.Println("P2:", P2)

	fmt.Println("--------------Step 6: Ranking--------------")
	fV1 := big.NewInt(1).Mul(v1, V1.X())
	fmt.Println("Player 1 final value:", fV1)
	fV2 := big.NewInt(1).Mul(v2, V2.X())
	fmt.Println("Player 2 final value:", fV2)

	sg := big.NewInt(0).Sub(fV1, fV2).Sign()
	if sg == 1 {
		if bankerVRF1.Verify(seed2.X(), V1, P1) {
			fmt.Println("Winner: Player 1")
		}
	} else {
		if bankerVRF2.Verify(seed2.X(), V2, P2) {
			fmt.Println("Winner: Player 2")
		}
	}

}
