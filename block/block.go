package block

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Timestamp int64
	Data []byte
	PrevBlockHash []byte
	Hash []byte
}

func (block *Block) ComputeHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

func New(data string, prevBlockHash []byte) *Block {
	newBlock := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	newBlock.ComputeHash()
	return newBlock
}
