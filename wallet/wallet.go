package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
			"log"
	"crypto/sha256"
	"github.com/itchyny/base58-go"
	"golang.org/x/crypto/ripemd160"
)

//basic wallet versions
const defaultWalletVersion = byte(0x01)
//default wallet file
const walletFile = "wallet.dat"
const addressChecksumLen = 4


type Wallet struct {
	Version byte
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func New() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{defaultWalletVersion, private, public}

	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatalf("failed generating private key for the wallet: %v\n", err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

func (w *Wallet) GetAddress() []byte {
	encoding := base58.BitcoinEncoding // or RippleEncoding or BitcoinEncoding
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{defaultWalletVersion}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address, err := encoding.Encode(fullPayload)
	if err != nil{
		log.Fatalf("Error encoding the new address: %v\n",err)
	}
	return address
}



func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil{
		log.Fatalf("Error encoding the new pubkey: %v\n",err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}