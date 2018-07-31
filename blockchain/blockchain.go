package blockchain
import "ZmeyCoin/block"

type Blockchain struct {
	blocks []*block.Block
	blocksCount uint64
}

func (blockchain *Blockchain) AddBlock(data string) {
	prevBlock := blockchain.blocks[blockchain.blocksCount - 1]
	newBlock := block.New(data, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
	blockchain.blocksCount++
}