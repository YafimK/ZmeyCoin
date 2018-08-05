// Loslly following https://en.bitcoin.it/wiki/Transaction

package transaction

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/ecdsa"
	"encoding/hex"
	"crypto/rand"
	"ZmeyCoin/util"
	"crypto/sha256"
	"crypto/elliptic"
	"math/big"
)

// The default reward for our loyal miner
//TODO: add some algorithm to change this as time passes by :D
const minerReward = 50

//Each transaction is constructed of several inputs and outputs
type Transaction struct {
	InCounter          int
	OutCounter         int
	TransactionInputs  []TransactionInput
	TransactionOutputs []TransactionOutput
}

func New(inputs []TransactionInput, outputs []TransactionOutput) *Transaction {
	//transaction.SetID()
	return &Transaction{InCounter: len(inputs), OutCounter: len(outputs), TransactionInputs: inputs, TransactionOutputs: outputs}
}

func (transaction *Transaction) String() string {
	return fmt.Sprintf("Transaction: \n %v inputs:\n %v \n %v outputs: %v", transaction.InCounter, transaction.OutCounter, transaction.TransactionInputs, transaction.TransactionOutputs)
}

func (transaction *Transaction) ToBytes() []byte {
	var container bytes.Buffer
	enc := gob.NewEncoder(&container) // Will write to network.
	err := enc.Encode(transaction)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return container.Bytes()
}

//Simple coinbase (first block transaction) generation transaction generator with no regard to script
func NewCoinbaseTransaction() *Transaction {
	txIn := TransactionInput{[]byte{}, -1, []byte{}, []byte{}}
	txOut := TransactionOutput{minerReward, []byte("")}
	transaction := New([]TransactionInput{txIn}, []TransactionOutput{txOut})

	return transaction
}

func (transaction *Transaction) IsCoinbase() bool {
	return transaction.InCounter == 1 && transaction.TransactionInputs[0].PrevTxOutIndex == -1
}

func (transaction *Transaction) GetMinimisedTransaction() *Transaction{

	var inputs []TransactionInput
	for _, input := range transaction.TransactionInputs {
		inputs = append(inputs, TransactionInput{input.PrevTransactionHash, input.PrevTxOutIndex, nil,nil})
	}
	outputs := append(make([]TransactionOutput, 0, len(transaction.TransactionOutputs)), transaction.TransactionOutputs...)

	return New(inputs, outputs)
}

func (transaction *Transaction) GetHash() []byte{
	return sha256.Sum256(util.SerializeObject(transaction))[:]
}

func (transaction *Transaction) SignTransaction(privateKey ecdsa.PrivateKey, previousTransactions map[string]Transaction) {
	if transaction.IsCoinbase() {
		return
	}

	minimizedTransaction := transaction.GetMinimisedTransaction()

	for inID, transactionInput := range minimizedTransaction.TransactionInputs {
		previousTransaction := previousTransactions[hex.EncodeToString(transactionInput.PrevTransactionHash)]
		minimizedTransaction.TransactionInputs[inID].SenderPublicKeyHash = previousTransaction.TransactionOutputs[transactionInput.PrevTxOutIndex].RecipientPubKeyHash
		minimizedTransactionId := minimizedTransaction.GetHash()
		minimizedTransaction.TransactionInputs[inID].SenderPublicKeyHash = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, minimizedTransactionId)
		if err !=nil {
			log.Panicf("Error during signning transaction: %v \n", err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		transaction.TransactionInputs[inID].Signature = signature
	}

}

func (transaction *Transaction) Verify(prevTXs map[string]Transaction) bool {
	minimizedTransaction := transaction.GetMinimisedTransaction()
	curve := elliptic.P256()

	for inID, vin := range transaction.TransactionInputs {
		prevTx := prevTXs[hex.EncodeToString(vin.PrevTransactionHash)]
		minimizedTransaction.TransactionInputs[inID].Signature = nil
		minimizedTransaction.TransactionInputs[inID].SenderPublicKeyHash = prevTx.TransactionOutputs[vin.PrevTxOutIndex].RecipientPubKeyHash
		minimizedTransactionHash := minimizedTransaction.GetHash()
		minimizedTransaction.TransactionInputs[inID].SenderPublicKeyHash = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.SenderPublicKeyHash)
		x.SetBytes(vin.SenderPublicKeyHash[:(keyLen / 2)])
		y.SetBytes(vin.SenderPublicKeyHash[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, minimizedTransactionHash, &r, &s) == false {
			return false
		}
	}

	return true
}