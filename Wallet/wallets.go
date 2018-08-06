package Wallet

import (
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"log"
	"io/ioutil"
	"os"
	)

//default Wallet file
const walletFile = "wallets.dat"
//TODO: change this to use db :)

type Wallets struct {
	Wallets map[string]*Wallet
}

func (wallets *Wallets) SaveToDisk() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (wallets *Wallets) LoadFromDisk() error{
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var walletsFromDisk Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&walletsFromDisk)
	if err != nil {
		log.Panic(err)
	}
	wallets.Wallets = walletsFromDisk.Wallets
	return nil
}

func (wallets *Wallets) GetWalletByAddress(address string) *Wallet{
	return wallets.Wallets[address]
}

func (wallets *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (wallets *Wallets) GetNewWallet() string{
	wallet := NewWallet()
	wallets.Wallets[wallet.GetWalletAddressString()] = wallet
	return wallet.GetWalletAddressString()
}

func New() *Wallets{
	wallets := &Wallets{make(map[string]*Wallet)}
	err := wallets.LoadFromDisk()
	if err != nil{
		log.Println("wallets file not found, will be created in the end of the next session or on manual save command")
	}
	//TODO: check if wallets file exist - then we should load it from disk else create one
	return wallets

	}