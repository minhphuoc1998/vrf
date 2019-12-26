package vrf

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func HCode(data string) (hcode string) {
	bits := make([]int, 0)

	for pos, i := 0, 0; i < len(data); pos++ {
		if isPerfectSquare(pos + 1) {
			bits = append(bits, 0)
		} else {
			bits = append(bits, cToB(data[i]))
			i++
		}
	}

	for pos, _ := range bits {
		if isPerfectSquare(pos + 1) {
			p := 0
			for i, _ := range bits {
				// Checks if the bit should be calculated
				if i+1 != pos+1 && ((i+1)&(pos+1) != 0) {
					p ^= bits[i]
				}
			}

			bits[pos] = p
		}
	}

	for _, bit := range bits {
		hcode += string(bit + 48)
	}
	return

}

func cToB(b byte) int {
	if b == '1' {
		return 1
	}
	return 0
}

func isPerfectSquare(n int) bool {
	return n == (n & -n)
}

func errorPosition(p []int) int {
	str := ""
	for _, val := range p {
		str = string(val+48) + str
	}

	number, _ := strconv.ParseInt(str, 2, 0)

	return int(number)
}

func BigToBin(x *big.Int) (s string) {
	return fmt.Sprintf("%b", x)
}

func PadLeft(x string, n int) (s string) {
	var str strings.Builder
	for i := 0; i < n - len(x); i ++ {
		str.WriteString("0")
	}
	str.WriteString(x)
	return str.String()
}