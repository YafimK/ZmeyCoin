package Interface

import (
	Interface2 "ZmeyCoin/Transaction/Interface"
)

type Chainstate interface {
	FindUnspentOutputsPerPubKeyHash(pubKeyHash []byte) []Interface2.TransactionOutput
	CountTransactions() int
	Reindex()
	FindSpendableOutputs(publicKeyHash []byte, amount int) (int, map[string][]int)
	Update(block interface{})
}