package Interface

import "ZmeyCoin/Transaction/Interface"

type Block interface {
ComputeHash()
ComputeTransactionsHash() []byte
String() string
Serialize() []byte
NewGenesisBlock(coinbaseTransaction *Interface.Transaction) Block
GetTransactions() []Interface.Transaction
}

