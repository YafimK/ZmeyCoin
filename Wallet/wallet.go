package Wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"ZmeyCoin/Util"
)

//basic Wallet versions
const defaultWalletVersion = byte(0x01)



//Default Wallet address checksum length
const addressChecksumLen = 4

type Wallet struct {
	Version    byte
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{defaultWalletVersion, private, public}

	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatalf("failed generating privateKey key for the Wallet: %v\n", err)
	}
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}

//Creates a new public Wallet address
func (w *Wallet) GetNewWalletAddress() []byte {
	pubKeyHash := Util.HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{defaultWalletVersion}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address, err := Util.EncodeInBase58(fullPayload)
	if err != nil {
		log.Fatalf("Error base58Encoder the new address: %v\n", err)
	}
	return address
}

func (w *Wallet) GetWalletAddressString() string {
	return string(w.GetNewWalletAddress())
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
