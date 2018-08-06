package Block

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
	"fmt"
	"ZmeyCoin/Transaction"
	"encoding/gob"
	"log"
	"ZmeyCoin/MerkleTree"
	"ZmeyCoin/util"
)

type Block struct {
	Timestamp int64
	Transactions []*Transaction.Transaction
	PrevBlockHash []byte
	Hash *[]byte
	Nonce int
}

func (block *Block) ComputeHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.ComputeTransactionsHash(), timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = &hash[:]
}

func (block *Block) ComputeTransactionsHash() []byte {
	var transactionHashes [][]byte

	for _, tx := range block.Transactions {
		transactionHashes = append(transactionHashes, tx.ToBytes())
	}

	merkleTree := MerkleTree.NewMerkleTree(transactionHashes)

	return merkleTree.Root.Data
}

func NewBlock(transactions []*Transaction.Transaction, prevBlockHash []byte) *Block {
	newBlock := &Block{Timestamp: time.Now().Unix(), Transactions: transactions, PrevBlockHash: prevBlockHash, Hash: nil,Nonce: 0}
	proofOfWork := ProofOfWork{BlockTip: newBlock}
	newBlock.Nonce, newBlock.Hash = proofOfWork.CalculateProof()
	return newBlock
}

func (block *Block) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Block creation timestamp: %x\n",  time.Unix(block.Timestamp, 0)))
	buffer.WriteString(fmt.Sprintf("Hash: %x\n", block.Hash))
	buffer.WriteString(fmt.Sprintf("Prev. hash: %x\n", block.PrevBlockHash))
	buffer.WriteString(fmt.Sprintf("Data: %v\n", block.Transactions))

	return fmt.Sprintf("%v", buffer.String())
}

// DeserializeBlock deserialize a Block
func DeserializeBlock(serializedBlock []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(serializedBlock))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (block *Block) Serialize() []byte {
	return util.SerializeObject(block)
}

