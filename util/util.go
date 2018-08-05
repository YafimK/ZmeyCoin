package util

import (
		"github.com/itchyny/base58-go"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
	"bytes"
	"encoding/gob"
)

func EncodeInBase58(targetPayload []byte) ([]byte, error){
	base58Encoder := base58.BitcoinEncoding
	address, err := base58Encoder.Encode(targetPayload)

	return address, err
}

func DecodeFromBase58(targetPayload []byte) ([]byte, error){
	base58Encoder := base58.BitcoinEncoding
	address, err := base58Encoder.Decode(targetPayload)

	return address, err
}


func HashPubKey(pubKey []byte) []byte {
	sha256PublicKeyEncoded := sha256.Sum256(pubKey)
	RIPEMD160encoder := ripemd160.New()
	_, err := RIPEMD160encoder.Write(sha256PublicKeyEncoded[:])
	if err != nil {
		log.Fatalf("Error encoding the new pubkey: %v\n", err)
	}
	hashedPublicKey := RIPEMD160encoder.Sum(nil)

	return hashedPublicKey
}


func SerializeObject(targetPayload interface{}) []byte{
	var serializedObjectBuffer bytes.Buffer
	enc := gob.NewEncoder(&serializedObjectBuffer)
	err := enc.Encode(targetPayload)
	if err != nil {
		log.Panic(err)
	}

	return serializedObjectBuffer.Bytes()
}