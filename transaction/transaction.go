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
	inCounter          int
	outCounter         int
	transactionInputs  []TransactionInput
	transactionOutputs []TransactionOutput
}

func New(inputs []TransactionInput, outputs []TransactionOutput) *Transaction {
	//transaction.SetID()
	return &Transaction{inCounter: len(inputs), outCounter: len(outputs), transactionInputs: inputs, transactionOutputs: outputs}
}

func (transaction *Transaction) String() string {
	return fmt.Sprintf("Transaction: \n %v inputs:\n %v \n %v outputs: %v", transaction.inCounter, transaction.outCounter, transaction.transactionInputs, transaction.transactionOutputs)
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
	return transaction.inCounter == 1 && transaction.transactionInputs[0].prevTxOutIndex == -1
}

func (transaction *Transaction) GetMinimisedTransaction() *Transaction{

	var inputs []TransactionInput
	for _, input := range transaction.transactionInputs {
		inputs = append(inputs, TransactionInput{input.prevTransactionHash, input.prevTxOutIndex, nil,nil})
	}
	outputs := append(make([]TransactionOutput, 0, len(transaction.transactionOutputs)), transaction.transactionOutputs...)

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

	for inID, transactionInput := range minimizedTransaction.transactionInputs {
		previousTransaction := previousTransactions[hex.EncodeToString(transactionInput.prevTransactionHash)]
		minimizedTransaction.transactionInputs[inID].SenderPublicKeyHash = previousTransaction.transactionOutputs[transactionInput.prevTxOutIndex].RecipientPubKeyHash
		minimizedTransactionId := minimizedTransaction.GetHash()
		minimizedTransaction.transactionInputs[inID].SenderPublicKeyHash = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, minimizedTransactionId)
		if err !=nil {
			log.Panicf("Error during signning transaction: %v \n", err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		transaction.transactionInputs[inID].Signature = signature
	}

}

func (transaction *Transaction) Verify(prevTXs map[string]Transaction) bool {
	minimizedTransaction := transaction.GetMinimisedTransaction()
	curve := elliptic.P256()

	for inID, vin := range transaction.transactionInputs {
		prevTx := prevTXs[hex.EncodeToString(vin.prevTransactionHash)]
		minimizedTransaction.transactionInputs[inID].Signature = nil
		minimizedTransaction.transactionInputs[inID].SenderPublicKeyHash = prevTx.transactionOutputs[vin.prevTxOutIndex].RecipientPubKeyHash
		minimizedTransactionHash := minimizedTransaction.GetHash()
		minimizedTransaction.transactionInputs[inID].SenderPublicKeyHash = nil

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

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, minimizedTransactionHash, &r, &s) == false {
			return false
		}
	}

	return true
}