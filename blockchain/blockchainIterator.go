package blockchain

import (
	"ZmeyCoin/block"
	)

type BlockchainIterator struct {
	blockCursor *block.Block

}

func NewBlockchainIterator(blockchain *Blockchain) *BlockchainIterator {

	return &BlockchainIterator{blockchain.blocks[blockchain.blocksCount - 1]}
}
// Next returns next block starting from the tip
func (blockchainIterator *BlockchainIterator) Next() *block.Block {
	nextBlockHash := blockchainIterator.blockCursor.PrevBlockHash
	deserializeBlock := block.DeserializeBlock(nextBlockHash)
	return deserializeBlock
}