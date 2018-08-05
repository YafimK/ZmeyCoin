package transaction

import (
	"fmt"
	"ZmeyCoin/util"
	"bytes"
)

type TransactionInput struct {
	PrevTransactionHash []byte
	PrevTxOutIndex      int
	//txInScript          []byte
	//txInScriptLength    int
	SenderPublicKeyHash []byte
	Signature           []byte
}

func (txInput *TransactionInput) String() string {
	//return fmt.Sprintf("Input: \n " +
	//	"PrevTransactionHash: %v \n" +
	//	"PrevTxOutIndex: %v \n " +
	//	"txInScript: %v \n " +
	//	"txInScript: %v\n", txInput.PrevTransactionHash,  txInput.PrevTxOutIndex, txInput.txInScript, txInput.txInScriptLength)
	return fmt.Sprintf("Input: \n "+
		"PrevTransactionHash: %v \n"+
		"PrevTxOutIndex: %v \n "+
		"SenderPublicKeyHash: %v \n "+
		"Signature: %v\n", txInput.PrevTransactionHash, txInput.PrevTxOutIndex, txInput.SenderPublicKeyHash, txInput.Signature)
}

func (txInput *TransactionInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := util.HashPubKey(txInput.SenderPublicKeyHash)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}


