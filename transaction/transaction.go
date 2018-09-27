
package Transaction

import (
	"ZmeyCoin/Transaction/Interface"
	"ZmeyCoin/Util"
	"ZmeyCoin/Wallet"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

// The default reward for our loyal miner
//TODO: add some algorithm to change this as time passes by :D
const minerReward = 50

//Each transaction is constructed of several inputs and outputs
type ZmeyTransaction struct{
	InCounter          int
	OutCounter         int
	TransactionInputs  []Interface.TransactionInput
	TransactionOutputs []Interface.TransactionOutput
}

//Simple coinbase (first Block transaction) generation transaction generator with no regard to script
func (transaction ZmeyTransaction) NewCoinbaseTransaction() Interface.Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}, []byte{}}
	txOut := NewTransactionOutput(minerReward, "")

	return NewZmeyTransaction([]TxInput, []TxOut)
}

func (transaction ZmeyTransaction) GetTransactionInputs() []Interface.TransactionInput {
	return transaction.TransactionInputs
}

func (transaction ZmeyTransaction) GetTransactionOutputs() []Interface.TransactionOutput {
	return transaction.TransactionOutputs
}

func (transaction ZmeyTransaction) String() string {
	return fmt.Sprintf("ZmeyTransaction: \n %v inputs:\n %v \n %v outputs: %v", transaction.InCounter, transaction.OutCounter, transaction.TransactionInputs, transaction.TransactionOutputs)
}


// NewUTXOTransaction creates a new transaction
func (transaction ZmeyTransaction) NewUTXOTransaction(from, to string, amount int, chainstate interface{}) Interface.Transaction {
	var inputs []TxInput
	var outputs []TxOut

	wallets, err := Wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWalletByAddress(from)
	pubKeyHash := Util.HashPubKey(wallet.PublicKey)
	acc, validOutputs := chainstate.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTransactionOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTransactionOutput(acc-amount, from)) // a change
	}

	tx := NewZmeyTransaction(inputs, outputs)
	chainstate.SignTransaction(tx, wallet.PrivateKey)

	return tx}

func NewZmeyTransaction(inputs []Interface.TransactionInput, outputs []Interface.TransactionOutput) ZmeyTransaction {
	//transaction.SetID()
	return ZmeyTransaction{InCounter: len(inputs), OutCounter: len(outputs), TransactionInputs: inputs, TransactionOutputs: outputs}
}
func (transaction ZmeyTransaction) ToBytes() []byte {
	var container bytes.Buffer
	enc := gob.NewEncoder(&container) // Will write to network.
	err := enc.Encode(transaction)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return container.Bytes()
}

func (transaction ZmeyTransaction) IsCoinbase() bool {
	return transaction.InCounter == 1 && transaction.TransactionInputs[0].GetPrevTxOutIndex() == -1
}

func (transaction ZmeyTransaction) GetMinimisedTransaction() Interface.Transaction {

	var inputs []TxInput
	for _, input := range transaction.TransactionInputs {
		inputs = append(inputs, TxInput{input.GetPrevTransactionHash(), input.GetPrevTxOutIndex(), nil,nil})
	}
	outputs := append(make([]TxOut, 0, len(transaction.TransactionOutputs)), transaction.GetTransactionOutputs()...)

	return NewZmeyTransaction(inputs, outputs)
}

func (transaction ZmeyTransaction) GetHash() []byte{
	//Store this as data member and make lazy eval
	sum256 := sha256.Sum256(Util.SerializeObject(transaction))
	return sum256[:]
}

func (transaction ZmeyTransaction) SignTransaction(privateKey ecdsa.PrivateKey, previousTransactions map[string]Interface.Transaction) {
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

func (transaction ZmeyTransaction) Verify(prevTXs map[string]Interface.Transaction) bool {
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
