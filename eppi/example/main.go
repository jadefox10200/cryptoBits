//version 1.0 X3DH/AMTP
//20200626_12:10
//go version go1.10 windows/amd64

package main

import (
	"fmt"

	"github.com/jadefox10200/x3dh"
)

func main() {

	b := x3dh.NewUser("Bob")
	bBundle := b.GetPreKeyBundle()
	
	fmt.Println()
	a := x3dh.NewUser("Alice")	
	aBundle := a.GetPreKeyBundle()
	// a.SendPreKeyBundle(os.Stdout)

	fmt.Println()	

	ekp, ask := a.MakeSecretAliceBundle(bBundle, "AMTP")
	fmt.Println()

	bsk, err:= b.MakeSecretBobBundle(aBundle, ekp, "AMTP")
	if err != nil {fmt.Println("Failed to make sk: %s", err.Error()); return}

	fmt.Printf("Bob's sk: %x\n", bsk)
	fmt.Printf("Alice's sk: %x\n", ask)

}	
