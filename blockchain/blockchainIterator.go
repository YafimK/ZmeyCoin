package BlockChain

import (
	"ZmeyCoin/Block"
		"github.com/dgraph-io/badger"
	"log"
)

type BlockchainIterator struct {
	cursorBlockHash *[]byte
	db              *badger.DB

}

func NewBlockchainIterator(blockchain *Blockchain) *BlockchainIterator {

	return &BlockchainIterator{&blockchain.BlockTip, blockchain.blockDb}
}
// Next returns next Block starting from the tip
func (blockchainIterator *BlockchainIterator) Next() *Block.Block {
	var block *Block.Block

	err := blockchainIterator.db.View(func(Txn *badger.Txn) error {
		item, err := Txn.Get(*blockchainIterator.cursorBlockHash)
		if err != nil {
			return err
		}
		var encodedBlock []byte
		encodedBlock, err = item.Value()
		if err != nil {
			return err
		}
		block = Block.DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	blockchainIterator.cursorBlockHash = &block.PrevBlockHash

	return block
}