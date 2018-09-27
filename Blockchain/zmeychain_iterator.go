package Blockchain

import (
	"ZmeyCoin/BlockChain/Interface"
	"github.com/dgraph-io/badger"
	"log"
)

type ZmeyChainIterator struct {
	cursorBlockHash *[]byte
	db              *badger.DB
}

func NewBlockchainIterator(blockchain *ZmeyChain) *ZmeyChainIterator {
	return &ZmeyChainIterator{&blockchain.BlockTip, blockchain.BlockDb}
}


// Next returns next ZmeyCoinBlock starting from the tip
func (blockchainIterator *ZmeyChainIterator) Next() Interface.Block {
	var block *ZmeyCoinBlock

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
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	blockchainIterator.cursorBlockHash = &block.PrevBlockHash

	return block
}