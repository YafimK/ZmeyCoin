package Interface

import (
	"crypto/ecdsa"
	"github.com/dgraph-io/badger"
	"ZmeyCoin/Transaction/Interface"
	)

type Blockchain interface {
	AddBlock(transactions []*Interface.Transaction) Block
	MineBlock(transactions []*Interface.Transaction) Block
	initBlockDb() (*badger.DB, error)
	initChainstateDb() (*badger.DB, error)
	PrintBlockChain()
	AddTransaction()
	FindTransaction(targetTransactionHash []byte) (Interface.Transaction, error)
	VerifyTransaction(targetTransaction *Interface.Transaction) bool
	Iterator() BlockchainIterator
	SignTransaction(tx *Interface.Transaction, privateKey ecdsa.PrivateKey)
	FindUnspentTransactionOutputs() map[string]Interface.TransactionOutputs
	GetChainStateDb() *badger.DB
}
