package blockchain

import (
	"ZmeyCoin/block"
	"fmt"
	"ZmeyCoin/transaction"
)

type Blockchain struct {
	blocks []*block.Block
	transactions []*transaction.Transaction //Transaction pending to be "block'ed"
	blocksCount int
}

func (blockchain *Blockchain) AddBlock(transactions []*transaction.Transaction) {
	prevBlock := blockchain.blocks[blockchain.blocksCount - 1]
	newBlock := block.New(transactions, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
	blockchain.blocksCount++
}

func (blockchain *Blockchain) MineBlock() {
	//TODO: gather all possible transactions and create a new block

	blockchain.AddBlock([]*transaction.Transaction{})
}

//we need to init the blockchain with genesis block
func New() *Blockchain {
	newBlockchain := Blockchain{}
	newBlockchain.blocks = append(newBlockchain.blocks,
		block.New([]*transaction.Transaction{transaction.NewCoinbaseTransaction()}, []byte{}))
	newBlockchain.blocksCount++
	return &newBlockchain
}

func (blockchain *Blockchain) printBlockChain() {
	fmt.Println("*** Blockchain ***")
	for index, curBlock := range blockchain.blocks {
		fmt.Printf("%v block\n",index)
		fmt.Println(curBlock)
	}
}

func (blockchain *Blockchain) addTransaction() {

}
