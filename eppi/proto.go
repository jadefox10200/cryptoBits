//version 1.1 X3DH/AMTP
//20200707_22:20 - UNTESTED
//go version go1.10 windows/amd64

package x3dh

import (
	"fmt"
	"os"
	"io"
  "io/ioutil"
	"encoding/json"
	"crypto/rand"
	"crypto/sha512"

	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/curve25519"		
)


func NewUser(n string) *User {

	u := &User{}
	u.Name = n
	PrivateBundle := GeneratePreKeyBundle()
	u.Bundle = PrivateBundle

	return u
}

func (u *User) SaveRemotePublicKeys(usrPath string, pb PublicPreKeyBundle) error {

  //Create a new "idk" file.
  idkNamePath := fmt.Sprintf("%s/idkpublic", usrPath)
  idkf, err := os.Create(idkNamePath)
  if err != nil {return err}

  _, err = idkf.Write([]byte(fmt.Sprintf("%x", pb.IdentityKeyPublic)))
  idkf.Close()
  if err != nil {return err}

  //Create a new "spk" file.
  spkNamePath := fmt.Sprintf("%s/spkpublic", usrPath)
  spkf, err := os.Create(spkNamePath)
  if err != nil {return err}

  spkf.Write([]byte(fmt.Sprintf("%x", pb.SignedPreKeyPublic)))
  spkf.Close()
  if err != nil {return err}

  //Create a new "sig" file.
  sigNamePath := fmt.Sprintf("%s/sig", usrPath)
  sigf, err := os.Create(sigNamePath)
  if err != nil { return err}

  sigf.Write([]byte(fmt.Sprintf("%x", pb.SignedPreKeySignature)))
  sigf.Close()
  if err != nil {return err}

  return nil

}

func (u *User) SavePublicKeys(usrPath string) error {

  //Create a new "idk" file.
  idkNamePath := fmt.Sprintf("%s/idkpublic", usrPath)
  idkf, err := os.Create(idkNamePath)
  if err != nil {return err}

  _, err = idkf.Write([]byte(fmt.Sprintf("%x", u.Bundle.IdentityKeyPublic)))
  idkf.Close()
  if err != nil {return err}

  //Create a new "spk" file.
  spkNamePath := fmt.Sprintf("%s/spkpublic", usrPath)
  spkf, err := os.Create(spkNamePath)
  if err != nil {return err}

  spkf.Write([]byte(fmt.Sprintf("%x", u.Bundle.SignedPreKeyPublic)))
  spkf.Close()
  if err != nil {return err}

  //Create a new "sig" file.
  sigNamePath := fmt.Sprintf("%s/sig", usrPath)
  sigf, err := os.Create(sigNamePath)
  if err != nil { return err}

  sigf.Write([]byte(fmt.Sprintf("%x", u.Bundle.SignedPreKeySignature)))
  sigf.Close()
  if err != nil {return err}

  return nil

}

func (u *User) LoadUserPrivateKeys(usrPath string) error {

  //Create a new "idk" file.
  idkNamePath := fmt.Sprintf("%s/idkprivate", usrPath)
  idkf, err := os.Open(idkNamePath)
  if err != nil {return err}

  idkbs, err := ioutil.ReadAll(idkf)
  if err != nil {return err}
  defer idkf.Close()

  //load private key into memory for use:
  u.Bundle.IdentityKeyPrivate = idkbs

  //Create a new "spk" file.
  spkNamePath := fmt.Sprintf("%s/spkprivate", usrPath)
  spkf, err := os.Open(spkNamePath)
  if err != nil {return err}

  spkfbs, err := ioutil.ReadAll(spkf)
  if err != nil {return err}

  u.Bundle.SignedPreKeyPrivate = spkfbs  
   
  return nil

}

func (u *User) LoadRemotePublicKeys(usrPath string) (PublicPreKeyBundle, error) {

  var pb = PublicPreKeyBundle{}

  idkNamePath := fmt.Sprintf("%s/idkpublic", usrPath)
  idkf, err := os.Open(idkNamePath)
  if err != nil {return PublicPreKeyBundle{}, err}

  idkbs, err := ioutil.ReadAll(idkf)
  defer idkf.Close()
  if err != nil {return PublicPreKeyBundle{}, err}
  
  pb.IdentityKeyPublic = idkbs  
  
  spkNamePath := fmt.Sprintf("%s/spkpublic", usrPath)
  spkf, err := os.Open(spkNamePath)
  if err != nil {return PublicPreKeyBundle{}, err}

  spkfbs, err := ioutil.ReadAll(spkf)
  defer spkf.Close()
  if err != nil {return PublicPreKeyBundle{}, err}
  
  pb.SignedPreKeyPublic = spkfbs    
  
  sigNamePath := fmt.Sprintf("%s/sig", usrPath)
  sigf, err := os.Open(sigNamePath)
  if err != nil {return PublicPreKeyBundle{}, err}

  sigbs, err := ioutil.ReadAll(sigf)
  defer spkf.Close()
  if err != nil {return PublicPreKeyBundle{}, err}

  pb.SignedPreKeySignature = sigbs

  return pb, nil

}

// type PublicPreKeyBundle struct {    

//     IdentityKeyPublic   []byte  //ipk    

//     SignedPreKeyID        uint32
//     SignedPreKeyPublic      []byte  //spk
//     SignedPreKeySignature   []byte

//     OneTimePreKeyID     int32
//     OneTimePreKeyPublic   []byte  //opk
// }

func (u *User) LoadUserPublicKeys(usrPath string) error {
  
  idkNamePath := fmt.Sprintf("%s/idkpublic", usrPath)
  idkf, err := os.Open(idkNamePath)
  if err != nil {return err}

  idkbs, err := ioutil.ReadAll(idkf)
  defer idkf.Close()
  if err != nil {return err}
  
  u.Bundle.IdentityKeyPublic = idkbs  
  
  spkNamePath := fmt.Sprintf("%s/spkpublic", usrPath)
  spkf, err := os.Open(spkNamePath)
  if err != nil {return err}

  spkfbs, err := ioutil.ReadAll(spkf)
  defer spkf.Close()
  if err != nil {return err}
  
  u.Bundle.SignedPreKeyPublic = spkfbs    
  
  sigNamePath := fmt.Sprintf("%s/sig", usrPath)
  sigf, err := os.Open(sigNamePath)
  if err != nil { return err}

  sigbs, err := ioutil.ReadAll(sigf)
  defer spkf.Close()
  if err != nil {return err}

  u.Bundle.SignedPreKeySignature = sigbs
  
  return nil

}

// type FullPreKeyBundle struct {    

//     IdentityKeyPublic   []byte
//     IdentityKeyPrivate    []byte

//     SignedPreKeyID          uint32
//     SignedPreKeyPublic      []byte
//     SignedPreKeyPrivate     []byte
//     SignedPreKeySignature   []byte

//     OneTimePreKeyID     int32
//     OneTimePreKeyPublic   []byte
//     OneTimePreKeyPrivate  []byte
// }


func (u *User) SavePrivateKeys(usrPath string) error {

  //Create a new "idk" file.
  idkNamePath := fmt.Sprintf("%s/idk", usrPath)
  idkf, err := os.Create(idkNamePath)
  if err != nil {return err}

  _, err = idkf.Write([]byte(fmt.Sprintf("%x", u.Bundle.IdentityKeyPrivate)))
  idkf.Close()
  if err != nil {return err}

  //Create a new "spk" file.
  spkNamePath := fmt.Sprintf("%s/spk", usrPath)
  spkf, err := os.Create(spkNamePath)
  if err != nil {return err}

  spkf.Write([]byte(fmt.Sprintf("%x", u.Bundle.SignedPreKeyPrivate)))
  spkf.Close()
  if err != nil {return err}

  return nil

}

func (u *User) SendPreKeyBundle(w io.Writer) {

	b := u.GetPreKeyBundle() 

	bs, err := json.Marshal(b)

	_, err = w.Write(bs)
	if err != nil {fmt.Println("Failed to write prekeybundle during send", err.Error()); return}

	return
}

func (u *User) GetPreKeyBundle() PublicPreKeyBundle {

	bundle := PublicPreKeyBundle{
		IdentityKeyPublic:		u.Bundle.IdentityKeyPublic,

	    SignedPreKeyID:       	u.Bundle.SignedPreKeyID,
	    SignedPreKeyPublic:     u.Bundle.SignedPreKeyPublic,
	    SignedPreKeySignature: 	u.Bundle.SignedPreKeySignature,

	    OneTimePreKeyID:  		u.Bundle.OneTimePreKeyID,
	    OneTimePreKeyPublic: 	u.Bundle.OneTimePreKeyPublic,
	}

	return bundle
}

func (u *User) VerifyIdentity(pb PublicPreKeyBundle) error {

	fmt.Println("Verified")

	var idpk [32]byte
	var sig [64]byte

	for i := 0; i < len(idpk); i++ {
		idpk[i] = pb.IdentityKeyPublic[i]		
	}

	for i := 0; i < len(sig); i++{
		sig[i] = pb.SignedPreKeySignature[i]
	}
		

	b := Verify( idpk , pb.SignedPreKeyPublic[:], &sig) 
	if !b {fmt.Println("failed to verify signature"); return fmt.Errorf("Vailed to verify signature")}

	return nil
}

func DH(remoteKey []byte, localKey []byte) ([]byte) {

  	var sk [32]byte 
  	var privateKey [32]byte 
  	var publicKey [32]byte  

  	if len(localKey) != 32 {fmt.Println("local key doesn't == 32"); os.Exit(1)}
  	if len(remoteKey) != 32 {fmt.Println("remote key doesn't == 32"); os.Exit(1)}

  	copy(privateKey[:], localKey[:32])
  	copy(publicKey[:], remoteKey[:32])
  	
  	curve25519.ScalarMult(&sk, &privateKey, &publicKey)
  
  	return sk[:]

 }

func (u *User) MakeSecretAliceBundle(pb PublicPreKeyBundle, info string) ([]byte, []byte) {

	u.VerifyIdentity(pb)

	ekPublic, ekPrivate := GenerateEphemeralKeyPair()

	dh1 := DH(pb.SignedPreKeyPublic, u.Bundle.IdentityKeyPrivate)
	dh2 := DH(pb.IdentityKeyPublic, ekPrivate)
	dh3 := DH(pb.SignedPreKeyPublic, ekPrivate)

	// dh1, _ := doubleratchet.DefaultCrypto{}.DH(ikA, spkB)
	// dh2, _ := doubleratchet.DefaultCrypto{}.DH(ekA, ikB)
	// dh3, _ := doubleratchet.DefaultCrypto{}.DH(ekA, spkB)

	fmt.Printf("DH 1 for %s is : %x\n", u.Name, dh1)
	fmt.Printf("DH 2 for %s is : %x\n", u.Name, dh2)
	fmt.Printf("DH 3 for %s is : %x\n", u.Name, dh3)

	//APPEND DH1-3 into key material for KDF. 
	km := append(dh1[:], dh2[:]...)
  	km = append(km, dh3[:]...)

  	//FOR NOW, DO NOT INCLUDE OPK
  	//if b.OneTimePreKey() != nil {
  //
  //	//	// fourth step with our ephemeral key
  //	//	// and the remote one time pre key
  //	//	dh4 := x.curve.KeyExchange(DHPair{
  //	//		PrivateKey: ephemeralKey.PrivateKey,
  //	//		PublicKey:  *b.OneTimePreKey(),
  //	//	})
  //
  //	//	km = append(km, dh4[:]...)
  //
  	//}

	sk, err := KDF(km, info)
	if err != nil {fmt.Printf("Failed to generate sk for %s: %s\n", u.Name, err.Error()); os.Exit(1)}

	return ekPublic, sk[:]
}


func (u *User) MakeSecretBobBundle(pb PublicPreKeyBundle, epk []byte, info string) ([]byte, error) {

	u.VerifyIdentity(pb)

	dh1 := DH(pb.IdentityKeyPublic, u.Bundle.SignedPreKeyPrivate)
	dh2 := DH(epk , u.Bundle.IdentityKeyPrivate)
	dh3 := DH(epk, u.Bundle.SignedPreKeyPrivate)

	fmt.Printf("DH 1 for %s is : %x\n", u.Name, dh1)
	fmt.Printf("DH 2 for %s is : %x\n", u.Name, dh2)
	fmt.Printf("DH 3 for %s is : %x\n", u.Name, dh3)

	//APPEND DH1-3 into key material for KDF. 
	km := append(dh1[:], dh2[:]...)
  	km = append(km, dh3[:]...)

	sk, err := KDF(km, info)
	if err != nil {fmt.Printf("Failed to generate sk for %s: %s\n", u.Name, err.Error()); return []byte{}, err}

	return sk[:], nil
}

func GenerateEphemeralKeyPair() ([]byte, []byte) {

	ekpPublic, ekpPrivate, err := GenerateKeyPair()
	if err != nil {fmt.Println("Failed to generate epk", err.Error()); os.Exit(1)}	

	return ekpPublic, ekpPrivate
}

func GenerateKeyPair() ([]byte, []byte, error){

	var privateKey [32]byte
  	if _, err := io.ReadFull(rand.Reader, privateKey[:]); err != nil {
  		return nil, nil, err
  	}
  	
  	privateKey[0] &= 248
  	privateKey[31] &= 127
  	privateKey[31] |= 64
  
  	var publicKey [32]byte
  	curve25519.ScalarBaseMult(&publicKey, &privateKey)

  	return publicKey[:], privateKey[:], nil
}

func GeneratePreKeyBundle() (FullPreKeyBundle) {

	var fullBundle = FullPreKeyBundle{}

	//GENERATE IDENTITY KEY PAIR: 	

	//Genereate the keys for signing: USE THIS FOR IDENTITY KEY PAIR
	// func GenerateKey(rand io.Reader) (publicKey PublicKey, privateKey PrivateKey, err error) {
	idkPublic, idkPrivate, err := GenerateKeyPair()
	if err != nil {fmt.Println("failed to generatekey", err.Error()); os.Exit(1)}

	fullBundle.IdentityKeyPublic = idkPublic
	fullBundle.IdentityKeyPrivate = idkPrivate

	//GENERATE SIGNED KEY PAIR:
	spkPublic, spkPrivate, signature := GenerateSignedPreKey(idkPrivate)	
	if err != nil {fmt.Println("Failed to generate spk", err.Error()); os.Exit(1)}	

	fullBundle.SignedPreKeyID = 1
	fullBundle.SignedPreKeyPrivate = spkPrivate
	fullBundle.SignedPreKeySignature = signature
	fullBundle.SignedPreKeyPublic = spkPublic

	opkPublic, opkPrivate, err := GenerateKeyPair()
	if err != nil {fmt.Println("Failed to generate opk", err.Error()); os.Exit(1)}	

	fullBundle.OneTimePreKeyID = 1
	fullBundle.OneTimePreKeyPrivate = opkPrivate
	fullBundle.OneTimePreKeyPublic = opkPublic

	return fullBundle
}

func PreFix() []byte {
  
	return []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

}

func KDF(keyMaterial []byte, info string) ([32]byte, error) {
  
	// create reader
	//func New(hash func() hash.Hash, secret, salt, info []byte) io.Reader
	r := hkdf.New(sha3.New256, append(PreFix(), keyMaterial...), nil, []byte(info))

	// fill the shared secret
	var secret [32]byte
	_, err := r.Read(secret[:])
	return secret, err

}

type User struct {
	Name			string	
	Bundle			FullPreKeyBundle	
}

type PublicPreKeyBundle struct {    

    IdentityKeyPublic		[]byte	//ipk    

    SignedPreKeyID       	uint32
    SignedPreKeyPublic      []byte	//spk
    SignedPreKeySignature 	[]byte

    OneTimePreKeyID  		int32
    OneTimePreKeyPublic 	[]byte	//opk
}

type FullPreKeyBundle struct {    

    IdentityKeyPublic		[]byte
    IdentityKeyPrivate		[]byte

    SignedPreKeyID        	uint32
    SignedPreKeyPublic     	[]byte
    SignedPreKeyPrivate    	[]byte
    SignedPreKeySignature 	[]byte

    OneTimePreKeyID  		int32
    OneTimePreKeyPublic 	[]byte
    OneTimePreKeyPrivate	[]byte
}

func GenerateSignedPreKey(private []byte) ([]byte, []byte, []byte) {
  	mPublic, mPrivate, err := GenerateKeyPair()
  	if err != nil {fmt.Println("Failed to generate keypair for signing", err.Error()); os.Exit(1)}
  
  	var random [64]byte
  	if _, err := io.ReadFull(rand.Reader, random[:]); err != nil {
  		panic(err)
  	}
  
  	var in [32]byte 

  	for k, _ := range in {
  		in[k] = private[k]
  	}

  	sig := Sign(&in, mPublic[:], random)

  	return	mPublic[:], mPrivate[:], sig[:]
  	
 }

func Sign(privateKey *[32]byte, message []byte, random [64]byte) *[64]byte {

  	// Calculate Ed25519 public key from Curve25519 private key
  	var A ExtendedGroupElement
  	var publicKey [32]byte
  	GeScalarMultBase(&A, privateKey)
  	A.ToBytes(&publicKey)
  
  	// Calculate r
  	diversifier := [32]byte{
  		0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
  		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
  		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
  		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
  
  	var r [64]byte
  	h := sha512.New()
  	h.Write(diversifier[:])
  	h.Write(privateKey[:])
  	h.Write(message)
  	h.Write(random[:])
  	h.Sum(r[:0])
  
  	// Calculate R
  	var rReduced [32]byte
  	ScReduce(&rReduced, &r)
  	var R ExtendedGroupElement
  	GeScalarMultBase(&R, &rReduced)
  
  	var encodedR [32]byte
  	R.ToBytes(&encodedR)
  
  	// Calculate S = r + SHA2-512(R || A_ed || msg) * a  (mod L)
  	var hramDigest [64]byte
  	h.Reset()
  	h.Write(encodedR[:])
  	h.Write(publicKey[:])
  	h.Write(message)
  	h.Sum(hramDigest[:0])
  	var hramDigestReduced [32]byte
  	ScReduce(&hramDigestReduced, &hramDigest)
  
  	var s [32]byte
  	ScMulAdd(&s, &hramDigestReduced, privateKey, &rReduced)
  
  	signature := new([64]byte)
  	copy(signature[:], encodedR[:])
  	copy(signature[32:], s[:])
  	signature[63] |= publicKey[31] & 0x80
  
  	return signature
}
  
  // Verify checks whether the message has a valid signature.
  func Verify(publicKey [32]byte, message []byte, signature *[64]byte) bool {
  
  	publicKey[31] &= 0x7F
  
  	/* Convert the Curve25519 public key into an Ed25519 public key.  In
  	particular, convert Curve25519's "montgomery" x-coordinate into an
  	Ed25519 "edwards" y-coordinate:
  	ed_y = (mont_x - 1) / (mont_x + 1)
  	NOTE: mont_x=-1 is converted to ed_y=0 since fe_invert is mod-exp
  	Then move the sign bit into the pubkey from the signature.
  	*/
  
  	var edY, one, montX, montXMinusOne, montXPlusOne FieldElement
  	FeFromBytes(&montX, &publicKey)
  	FeOne(&one)
  	FeSub(&montXMinusOne, &montX, &one)
  	FeAdd(&montXPlusOne, &montX, &one)
  	FeInvert(&montXPlusOne, &montXPlusOne)
  	FeMul(&edY, &montXMinusOne, &montXPlusOne)
  
  	var A_ed [32]byte
  	FeToBytes(&A_ed, &edY)
  
  	A_ed[31] |= signature[63] & 0x80
  	signature[63] &= 0x7F
  
  	var sig = make([]byte, 64)
  	var aed = make([]byte, 32)

  	copy(sig, signature[:])
  	copy(aed, A_ed[:])

  	return ed25519.Verify(aed, message, sig)
  }
