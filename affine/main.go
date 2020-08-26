//This shows that despite doing a double decryption with the affine cypher, the key space is the same and therefore 
//it does not become more secure. All that changes is the value of the keys. It may "seem" more secure but in the end,
//an attacker wouldn't be computing a3. They would simply find a3 directly.

//go version go1.10 windows/amd64
package main

import (
	"fmt"
	"bufio"
	"strings"
	"math/big"
)

func main() {
	

	var a1 int64 = 3
	var b1 int64 = 5
	var a2 int64 = 11
	var b2 int64 = 7

	var m int64 = 26

	var a3 int64 = (a1 * a2) % m
	var b3 int64 = ((a2 * b1) + b2) % m

	// var a3 int64 = 7
	// var b3 int64 = 10

	x := "thisisclear"

	y1 := encrypt(a1, b1, m, x)

	fmt.Printf("%v <- Single encryption", y1)

	y2 := encrypt(a2, b2, m, y1)

	fmt.Printf("%v <- Double encryption", y2)

	y3 := encrypt(a3, b3, m, x)

	fmt.Println("%v <- Simulated double encryption", y3)

	x1 := decrypt(a1, b1, m, y1)

	fmt.Println("%v <- Single decryption", x1)

	x3 := decrypt(a3, b3, m, y2)

	fmt.Println("%v <- Decryption of double encryption using computed keys", x3)


	return
}

func decrypt(a int64, b int64, m int64, enc string) string {

	var txt strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(enc))

	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {

		//a-1 * (y-b)
		
		x := big.NewInt(0).Mod(big.NewInt(0).Mul(big.NewInt(0).ModInverse(big.NewInt(a), big.NewInt(m)), big.NewInt(0).Sub( big.NewInt(int64(scanner.Bytes()[0] - 97 )), big.NewInt(b) ) ), big.NewInt(m)).Int64() + 97
		txt.WriteRune(rune(x))

	}


	return txt.String()
}


func encrypt(a int64, b int64, m int64, text string) string {

	var enc strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(text))

	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {

		//a * x + b

		x := big.NewInt(0).Mod(	big.NewInt(0).Add(big.NewInt(0).Mul( big.NewInt(a), big.NewInt(int64(scanner.Bytes()[0] - 97) ) ) ,big.NewInt(b)  ), big.NewInt(m)).Int64() + 97


		enc.WriteRune(rune(x))

	}


	return enc.String()
}
