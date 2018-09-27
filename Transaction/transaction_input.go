package Transaction

import (
	"fmt"
	"ZmeyCoin/Util"
	"bytes"
)

type TxInput struct {
	PrevTransactionHash []byte
	PrevTxOutIndex      int
	//txInScript          []byte
	//txInScriptLength    int
	SenderPublicKeyHash []byte
	Signature           []byte
}

func (txInput *TxInput) String() string {
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

func (txInput *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := Util.HashPubKey(txInput.SenderPublicKeyHash)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (txInput *TxInput) GetPrevTransactionHash() []byte {
	return txInput.PrevTransactionHash
}
func (txInput *TxInput) GetPrevTxOutIndex() int {
	return txInput.PrevTxOutIndex
}
