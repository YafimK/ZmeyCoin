package blockchain

import (
	"ZmeyCoin/Block"
	)

type BlockchainIterator struct {
	blockCursor *Block.Block

}

func NewBlockchainIterator(blockchain *Blockchain) *BlockchainIterator {

	return &BlockchainIterator{blockchain.blocks[blockchain.blocksCount - 1]}
}
// Next returns next Block starting from the tip
func (blockchainIterator *BlockchainIterator) Next() *Block.Block {
	if blockchainIterator.blockCursor.PrevBlockHash == nil || len (blockchainIterator.blockCursor.PrevBlockHash) == 0{
		return nil
	}
	nextBlockHash := blockchainIterator.blockCursor.PrevBlockHash
	deserializeBlock := Block.DeserializeBlock(nextBlockHash)
	return deserializeBlock
}