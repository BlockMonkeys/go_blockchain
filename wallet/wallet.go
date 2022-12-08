package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"pkg/utils"
)

// Create Msg Hash

// Create Sign With Private Key

// Verify (Message + Signature) & Public Key => Verification
const (
	fileName string = "my.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string // Hex String
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644) // Read & Write Permission
	utils.HandleErr(err)
}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	// Restore Key
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	//합치기
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleErr(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	// Sig -> R & S
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	// Div Signature
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	// Create Public Sig
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)

	// Create Public Key
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	// Payload -> Bytes
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	//Verify
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		// 월렛존재함.
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
		// Restore From File
		// No => Create Private Key, Save to File
		w.Address = aFromK(w.privateKey)
	}
	return w
}
