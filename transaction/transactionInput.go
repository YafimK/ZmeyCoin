package transaction

import (
	"fmt"
	"ZmeyCoin/util"
	"bytes"
)

type TransactionInput struct {
	prevTransactionHash []byte
	prevTxOutIndex      int
	//txInScript          []byte
	//txInScriptLength    int
	SenderPublicKeyHash []byte
	Signature           []byte
}

func (txInput *TransactionInput) String() string {
	//return fmt.Sprintf("Input: \n " +
	//	"prevTransactionHash: %v \n" +
	//	"prevTxOutIndex: %v \n " +
	//	"txInScript: %v \n " +
	//	"txInScript: %v\n", txInput.prevTransactionHash,  txInput.prevTxOutIndex, txInput.txInScript, txInput.txInScriptLength)
	return fmt.Sprintf("Input: \n "+
		"prevTransactionHash: %v \n"+
		"prevTxOutIndex: %v \n "+
		"SenderPublicKeyHash: %v \n "+
		"Signature: %v\n", txInput.prevTransactionHash, txInput.prevTxOutIndex, txInput.SenderPublicKeyHash, txInput.Signature)
}

func (txInput *TransactionInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := util.HashPubKey(txInput.SenderPublicKeyHash)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}


