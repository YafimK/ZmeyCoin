package blockchain

import (
	"ZmeyCoin/block"
	"fmt"
	"ZmeyCoin/transaction"
)

type Blockchain struct {
	blocks []*block.Block
	blocksCount int
}

func (blockchain *Blockchain) AddBlock(data string) {
	prevBlock := blockchain.blocks[blockchain.blocksCount - 1]
	newBlock := block.New(data, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
	blockchain.blocksCount++
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
