//UNDERSTANDING CRYPTOGRAPHY
//Author Christof Parr
//pg 40

//This code was written to express a formula given on the above mentioned page. I am sure there are better ways to write this code, but it works. 

//PRNG reverse calculation of A and B where KEY(A, B)
//A and B are in theory unknown. However, if S1, S2 and S3 are known, A and B can be calculated. 
//In practice, the bits of S1, S2, S3... (aka Si) would be used to encrypt data for example in a stream cipher.
//Such that Si = (s1, s2, s3, ... m-1) where m would be the size of Si. Where xi represent the bits of clear text
// yi == xi + si mod 2. 
//If an attacker knows the first pieces of clear text which would comprise S1 - S3 the following formula would be done to compute S1 - S3. 
// si == yi + x mod m

//The below code simply shows that if you have calculated S1, S2 and S3 correctly using the above you can calculate A and B. 
// The variable gcd must = 1 (or decryption wouldn't be possible as there would be no inverse)
// The variable gcdInverse must also = 1 for the formula to work. When this is not the case, there can be more than one answer for A based on this formula. 
// Further code could be written to check the other possible answers for A. Therefore, A could still be calculated but would require checking the range of options available. 

//go version go1.10 windows/amd64
package main

import (
	"fmt"
	"math/big"
)

func main() {
	
	UsingBig()

	return
}

func UsingBig() {

fmt.Println("-----------------------------------------")
	z := big.NewInt(0)

	//ANSI C spec: 
	S0 := big.NewInt(12345).Int64()
	A := big.NewInt(1103515245)
	fmt.Printf("A: %v\n", A.Int64())
	B := big.NewInt(12345)
	fmt.Printf("B: %v\n", B.Int64())

	m := big.NewInt(2147483648).Int64()

	//Using small numbers to build the formula
	// S0 := big.NewInt(12).Int64()

	// A := big.NewInt(5)
	// fmt.Printf("A: %v\n", A.Int64())
	// B := big.NewInt(9)
	// fmt.Printf("B: %v\n", B.Int64())
	// m := big.NewInt(26).Int64()

	S1 := z.Mod(z.Add(z.Mul(A, big.NewInt(S0)), B), big.NewInt(m)).Int64()
	S2 := z.Mod(z.Add(z.Mul(A, big.NewInt(S1)), B), big.NewInt(m)).Int64()	
	S3 := z.Mod(z.Add(z.Mul(A, big.NewInt(S2)), B), big.NewInt(m)).Int64()


	//A == (S2 - S3) / (S1 - S2) mod m

	// (S2 - S3) <- This must be then inverse mul of m so that we can multiply it by a2
	a1 := z.ModInverse(z.Sub(big.NewInt(S2), big.NewInt(S3)), big.NewInt(m) ).Int64()
	
	// (S1 - S2)
	a2 := z.Sub(big.NewInt(S1), big.NewInt(S2)).Int64()

	// a = a1 * (s1 - s2) mod m
	Acomp := z.ModInverse( z.Mul( big.NewInt(a1), big.NewInt(a2)), big.NewInt(m) ).Int64()	

	fmt.Printf("Computed A: %v\n", Acomp)

	// (S2 - S3) <- This must be then inverse mul of m so that we can multiply it by b2
	//NOTE: b1 == a1
	b1 := z.ModInverse(z.Sub(big.NewInt(S2), big.NewInt(S3)), big.NewInt(m) ).Int64()

	// (S1 - S2) NOTE: b2 == a2
	b2 := z.Sub(big.NewInt(S1), big.NewInt(S2)).Int64()

	//b3 = (S2 - S3) / (S1 - S2) <- NOTE b3 == A
	b3 := z.ModInverse(z.Mul( big.NewInt(b1), big.NewInt(b2) ), big.NewInt(m)).Int64()

	b4 := z.Mul( big.NewInt(S1), big.NewInt(b3) ).Int64()
	b5 := z.Sub( big.NewInt(S2), big.NewInt(b4)).Int64()

	//B = S2 - S1 * (S2 - S3) / (S1 - S2)
	Bcomp := z.Mod( big.NewInt(b5), big.NewInt(m)).Int64()

	fmt.Printf("a1: %v\n",a1)
	fmt.Printf("a2: %v\n",a2)

	fmt.Printf("Computed B: %v\n", Bcomp)
	fmt.Printf("b1: %v\n",b1)
	fmt.Printf("b2: %v\n",b2)
	fmt.Printf("b3: %v\n",b3)
	fmt.Printf("b4: %v\n",b4)
	fmt.Printf("b5: %v\n",b5)	

	fmt.Printf("S1: %v\n",S1)
	fmt.Printf("S2: %v\n",S2)
	fmt.Printf("S3: %v\n",S3)

	fmt.Println("-----------------------------------------")

	 gcdInverse := big.NewInt(0).GCD(nil, nil, big.NewInt(int64(S1-S2)), big.NewInt(int64(m)))
	 gcd := big.NewInt(0).GCD(nil, nil, A, big.NewInt(int64(m)))

	 fmt.Printf("gcdInverse: %v\n", gcdInverse)
	 fmt.Printf("gcd: %v\n", gcd)

}
