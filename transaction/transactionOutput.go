package transaction

import (
	"fmt"
	"ZmeyCoin/util"
	"log"
	"bytes"
)

type TransactionOutput struct {
	value int
	//txOutScript       []byte
	//txOutScriptLength int
	RecipientPubKeyHash []byte
}

func (txOutput *TransactionOutput) String() string {
	//return fmt.Sprintf("Output: \n " +
	//	"value: %v \n" +
	//	"txOutScript: %v \n " +
	//	"txOutScriptLength: %v \n ",
	//	txOutput.value,  txOutput.txOutScript, txOutput.txOutScriptLength)
	return fmt.Sprintf("Output: \n "+
		"value: %v \n"+
		"RecipientPublicKeyHash: %v \n ", txOutput.value, txOutput.RecipientPubKeyHash)
}

//This decodes the public key from the address we've received and enables us to lock the output with the specific amount
func (txOutput *TransactionOutput) Lock(address []byte) {
	//TODO: add verify that you cannot call this method twice
	pubKeyHash, err := util.DecodeFromBase58(address)
	if err != nil {
		log.Fatalf("Error during base58 Encoding the new address: %v\n", err)
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	txOutput.RecipientPubKeyHash = pubKeyHash
}

func (txOutput *TransactionOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txOutput.RecipientPubKeyHash, pubKeyHash) == 0
}


