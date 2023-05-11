package main

import (
    signal "github.com/dosco/libsignal-protocol-go"
    "golang.org/x/crypto/curve25519"
    "fmt"
)

func main() {


    msA := signal.NewMemoryStore()
    msB := signal.NewMemoryStore()

    // Alice keys
    idA := signal.GenerateRegistrationId()
    msA.PutLocalRegistrationID(idA)


    ikpA := signal.GenerateIdentityKeyPair()
    msA.PutIdentityKeyPair(ikpA)

    pkA := signal.GeneratePreKey(1)
    msA.PutPreKey(1, pkA)

    spkA := signal.GenerateSignedPreKey(ikpA, 1)
    msA.PutSignedPreKey(1, spkA)

    // Bob keys

    idB := signal.GenerateRegistrationId()
    msB.PutLocalRegistrationID(idB)

    ikpB := signal.GenerateIdentityKeyPair()
    msB.PutIdentityKeyPair(ikpB)

    pkB := signal.GeneratePreKey(1)
    msB.PutPreKey(1, pkB)

    spkB := signal.GenerateSignedPreKey(ikpB, 1)
    msB.PutSignedPreKey(1, spkB)

    fmt.Printf("== Long-term keys\n")

    fmt.Printf("Alice ID (idA): %d\n", idA)
    fmt.Printf("Alice ID key (ikpA.Priv): %x\n", *ikpA.Priv)
    fmt.Printf("Alice ID key (ikpA.Pub): %x\n", *ikpA.Pub)
    fmt.Printf("Alice Pre-key (spkA.Priv): %x\n", *spkA.Priv)
    fmt.Printf("Alice Pre-key (spkA.Pub): %x\n\n", *spkA.Pub)

    fmt.Printf("Bob ID (idB): %d\n", idB)
    fmt.Printf("Bob ID key (ikpA.Priv): %x\n", *ikpB.Priv)
    fmt.Printf("Bob ID key (ikpA.Pub): %x\n", *ikpB.Pub)
    fmt.Printf("Bob Pre-key (spkA.Priv): %x\n", *spkB.Priv)
    fmt.Printf("Bob Pre-key (spkA.Pub): %x\n\n", *spkB.Pub)


    fmt.Printf("== Ephermal keys\n")
    ekA := signal.GenerateEphemeralKeyPair()

    fmt.Printf("Alice Ephemeral key (ekA.Priv): %x\n", *ekA.Priv)
    fmt.Printf("Alice Ephemeral key (ekA.Pub): %x\n", *ekA.Pub)


    // Calculate DH parameters and keys

    dh1 := signal.DH(ikpA.Priv, spkB.Pub)
    dh2 := signal.DH(ekA.Priv, ikpB.Pub)
    dh3 := signal.DH(ekA.Priv, spkB.Pub)


    fmt.Printf("\n== Alice calculates\n")

    fmt.Printf("DH1: %x\n", *dh1)
    fmt.Printf("DH2: %x\n", *dh2)
    fmt.Printf("DH3: %x\n", *dh3)

    dhList := [][]byte{dh1[:], dh2[:], dh3[:]}

    res := signal.KDF(dhList...)

    fmt.Printf("\nKey (RootKey, ChainKey, Index): %x\n\n", *res)


    dh1 = signal.DH(spkB.Priv, ikpA.Pub)
    dh2 = signal.DH(ikpB.Priv, ekA.Pub)
    dh3 = signal.DH(spkB.Priv, ekA.Pub)

    fmt.Printf("== Bob calculates\n")

    fmt.Printf("DH1: %x\n", *dh1)
    fmt.Printf("DH2: %x\n", *dh2)
    fmt.Printf("DH3: %x\n", *dh3)

    dhList = [][]byte{dh1[:], dh2[:], dh3[:]}

    res = signal.KDF(dhList...)

    fmt.Printf("\nKey (RootKey, ChainKey, Index): %x\n", *res)

    //SENDER IK GOES FIRST. RECEIVER IK GOES SECOND. THESE ARE CONCAT TOGETHER
    // AD = Encode(IKA) || Encode(IKB)
    var ad1 [32]byte
    var ad2 [32]byte

    curve25519.ScalarBaseMult(&ad1, &ikpA.Priv.Key)

    curve25519.ScalarBaseMult(&ad2, &ikpB.Pub.Key)

    ad := [][]byte{ad1[:], ad2[:]}

    fmt.Printf("\nad key: %x\n", ad)
    
}

// A sample run is defined below. Notice that Bob and Alice get the same values for DH1, DH2 and DH3, and can thus generate the same key:

// Alice ID (idA): 2574934372042371790
// Alice ID key (ikpA.Priv): {2017589787df9f2bf1a4acf827f3fb0c18817e8a7d30a7e7de6ca60de61b5f5e}
// Alice ID key (ikpA.Pub): {31cb48103c737e3d1ae52779183b1b0148bd9bc4d3980ddf71aad7f9b310d176}
// Alice Pre-key (spkA.Priv): {b0ba5ae7f54f53c4d124d05564ddd427d749535cc831dbb44b90a13cd9056246}

// Alice Pre-key (spkA.Pub): {8d4ffe9269188c0771c6a03e319ae63aab066066d130a94886b8b4cce3fe0366}

// Bob ID (idB): 2919514492885884887
// Bob ID key (ikpA.Priv): {e819e8c065392c9234b87360e14256fa71582b94cf4764bbcb955a6fd6e64b61}
// Bob ID key (ikpA.Pub): {ee0200274d4daa0f83bafb3757f09cb83095d28491b66ab7778dc8133a2cd105}
// Bob Pre-key (spkA.Priv): {f0930e2223fa35111253c63c010774e7db6e4007ed8df08907b386fff8dcb073}

// Bob Pre-key (spkA.Pub): {fcedada3c0f49dacef7a32a9e109b8f4e4fd8dd4c82bcc13a08d2b82efd70947}


// Alice Ephemeral key (ekA.Priv): {284446938ed2815ad7de21f675a763220e39e81d2bd1f964de78bde01041886c}
// Alice Ephemeral key (ekA.Pub): {353299ab7c1acb4a4bd349c3421ee91e025ac3c576ce308b8a1045f1c7f50411}

// == Alice calculates
// DH1: 9d0a39125ffd0ff12c7307e727b5191d1d9b65bd312b73e07915ba92b924d82b
// DH2: 357559d2e175285c3bfac101faa8523342946b720434abedc0aa0d1b41e56809
// DH3: b7843e2f1c98391fad3a2d8a00f0b0a3799468fedda1e5fc53dcb2a6a40b0d37

// Key (RootKey, ChainKey, Index): {39ca2aa1e04a5b8a85b0f4b61074bd6cd30fdcb856aba18f34a46ee87ea95652 d2845861ee9dc97606d44e1bbf32bb65bed79b838835d79be0b5cefc762dc254 0}


// == Bob calculates
// DH1: 9d0a39125ffd0ff12c7307e727b5191d1d9b65bd312b73e07915ba92b924d82b
// DH2: 357559d2e175285c3bfac101faa8523342946b720434abedc0aa0d1b41e56809
// DH3: b7843e2f1c98391fad3a2d8a00f0b0a3799468fedda1e5fc53dcb2a6a40b0d37

// Key (RootKey, ChainKey, Index): {39ca2aa1e04a5b8a85b0f4b61074bd6cd30fdcb856aba18f34a46ee87ea95652 d2845861ee9dc97606d44e1bbf32bb65bed79b838835d79be0b5cefc762dc254 0}
