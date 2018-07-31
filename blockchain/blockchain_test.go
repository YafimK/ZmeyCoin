package blockchain

import (
	"testing"
	"bytes"
)

func TestBlockchain_AddBlock(test *testing.T) {
	blockchain := New()
	blockchain.AddBlock("first block")
	blockchain.AddBlock("second block")
	blockchain.AddBlock("third block")

	if len(blockchain.blocks) != blockchain.blocksCount {
		test.Errorf("the block counter isn'test working right got %v but exepected %v",blockchain.blocksCount, len(blockchain.blocks))
	}

	if bytes.Equal(blockchain.blocks[2].Data, []byte("second block")) {
		test.Errorf("The data in the second block doesn't match the expected data that should be in it")
	}
}