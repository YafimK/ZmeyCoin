package block

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
	"fmt"
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

func (block *Block) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Block creation timestamp: %x\n",  time.Unix(block.Timestamp, 0)))
	buffer.WriteString(fmt.Sprintf("Hash: %x\n", block.Hash))
	buffer.WriteString(fmt.Sprintf("Prev. hash: %x\n", block.PrevBlockHash))
	buffer.WriteString(fmt.Sprintf("Data: %s\n", block.Data))

	return fmt.Sprintf("%v", buffer.String())
}