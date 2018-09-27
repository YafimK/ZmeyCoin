
package Interface

import (
	"crypto/ecdsa"
)

type Transaction interface {
String() string
ToBytes() []byte
NewCoinbaseTransaction() Transaction
IsCoinbase() bool
GetMinimisedTransaction() Transaction
GetHash() []byte
SignTransaction(privateKey ecdsa.PrivateKey, previousTransactions map[string]Transaction)
Verify(prevTXs map[string]Transaction) bool
NewUTXOTransaction(from, to string, amount int, chainstate interface{}) Transaction
GetTransactionInputs() []TransactionInput
GetTransactionOutputs() []TransactionOutput
}
