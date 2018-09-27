package Transaction

import (
	"ZmeyCoin/Transaction/Interface"
	"ZmeyCoin/Util"
	"bytes"
	"fmt"
	"log"
	"sync"
)

//TxOut
type TxOut struct {
	value int
	//txOutScript       []byte
	//txOutScriptLength int
	RecipientPubKeyHash []byte
	recipentHashLock sync.Once
}

func (txOutput *TxOut) GetValue() interface{} {
	return txOutput.value
}

func (txOutput *TxOut) NewTransactionOutput(amount int, recipient string) *Interface.TransactionOutput {
	panic("implement me")
}

func NewTransactionOutput(amount int,  recipient string) *TxOut {
	output := &TxOut{value: amount, RecipientPubKeyHash: nil}
	output.lockOutputByAddress([]byte(recipient))
	return output
}

func (txOutput *TxOut) String() string {
	//return fmt.Sprintf("Output: \n " +
	//	"value: %v \n" +
	//	"txOutScript: %v \n " +
	//	"txOutScriptLength: %v \n ",
	//	txOutput.value,  txOutput.txOutScript, txOutput.txOutScriptLength)
	return fmt.Sprintf("Output: \n "+
		"value: %v \n"+
		"RecipientPublicKeyHash: %v \n ", txOutput.value, txOutput.RecipientPubKeyHash)
}

//This decodes the public key from the address we've received and enables us to lockOutputByAddress the output with the specific amount
func (txOutput *TxOut) lockOutputByAddress(address []byte) {
	//You cannot lockOutputByAddress same output twice
	txOutput.recipentHashLock.Do(func(){
		pubKeyHash, err := Util.DecodeFromBase58(address)
		if err != nil {
			log.Fatalf("Error during base58 Encoding the new address: %v\n", err)
		}
		pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
		txOutput.RecipientPubKeyHash = pubKeyHash
	})
}

func (txOutput *TxOut) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txOutput.RecipientPubKeyHash, pubKeyHash) == 0
}


