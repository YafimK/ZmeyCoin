package ZmeyCoin

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

func (b* Block) ComputeHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func New(data string, prevBlockHash []byte) *Block {
	newBlock := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	newBlock.ComputeHash()
	return newBlock
}
