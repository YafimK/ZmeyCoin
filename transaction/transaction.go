// Loslly following https://en.bitcoin.it/wiki/Transaction

package transaction

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

// The default reward for our loyal miner
//TODO: add some algorithm to change this as time passes by :D
const minerReward = 50

//Each transaction is constructed of several inputs and outputs
type Transaction struct {
	inCounter  int
	outCounter int
	in         []TXInput
	out        []TXOutput
}

func New(inputs []TXInput, outputs []TXOutput) *Transaction{
	//transaction.SetID()
	return &Transaction{len(inputs), len(outputs), inputs, outputs}
}

func (tx *Transaction) String () string{
	return fmt.Sprintf("Transaction: \n %v inputs:\n %v \n %v outputs: %v", tx.inCounter,  tx.outCounter, tx.in, tx.out)
}

func (tx *Transaction) ToBytes () []byte {
	var container bytes.Buffer
	enc := gob.NewEncoder(&container) // Will write to network.
	err := enc.Encode(tx)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return container.Bytes()
}

type TXInput struct {
	prevTransactionHash []byte
	prevTxOutIndex      int
	txInScript          []byte
	txInScriptLength    int
}

func (txInput *TXInput) String () string{
	return fmt.Sprintf("Input: \n " +
		"prevTransactionHash: %v \n" +
		"prevTxOutIndex: %v \n " +
		"txInScript: %v \n " +
		"txInScript: %v\n", txInput.prevTransactionHash,  txInput.prevTxOutIndex, txInput.txInScript, txInput.txInScriptLength)
}

type TXOutput struct {
	value             int
	txOutScript       []byte
	txOutScriptLength int
}

func (txOutput *TXOutput) String () string{
	return fmt.Sprintf("Output: \n " +
		"value: %v \n" +
		"txOutScript: %v \n " +
		"txOutScriptLength: %v \n ",
		txOutput.value,  txOutput.txOutScript, txOutput.txOutScriptLength)
}

//Simple coinbase (first block transaction) generation transaction generator with no regard to script
func NewCoinbaseTransaction() *Transaction {
	txIn := TXInput{[]byte{}, -1, []byte(""),len("")}
	txOut := TXOutput{minerReward, []byte(""), len("")}
	transaction := New([]TXInput{txIn}, []TXOutput{txOut})

	return transaction
}

//TODO: create toString methods..
//We will use the most common bitcoin script when used to deliver coins to public key hash (P2PKH)
//Pubkey script: OP_DUP OP_HASH160 <PubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
//Signature script: <sig> <pubkey>
func (txInput *TXInput) LockInput() {

}