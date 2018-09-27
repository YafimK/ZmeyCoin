package Blockchain

import (
	Interface2 "ZmeyCoin/BlockChain/Interface"
	"ZmeyCoin/MerkleTree"
	"ZmeyCoin/Transaction/Interface"
	"ZmeyCoin/Util"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"time"
)

type ZmeyCoinBlock struct {
	Timestamp int64
	Transactions []Interface.Transaction
	PrevBlockHash []byte
	Hash *[]byte
	Nonce int
}

func (block ZmeyCoinBlock) GetTransactions() []Interface.Transaction {
	return block.Transactions
}

func NewBlock(transactions []Interface.Transaction, prevBlockHash []byte) ZmeyCoinBlock {
	newBlock := ZmeyCoinBlock{Timestamp: time.Now().Unix(), Transactions: transactions, PrevBlockHash: prevBlockHash, Hash: nil,Nonce: 0}
	proofOfWork := ProofOfWork{BlockTip: &newBlock}
	newBlock.Nonce, newBlock.Hash = proofOfWork.CalculateProof()
	return newBlock
}

func DeserializeBlock(serializedBlock []byte) *ZmeyCoinBlock {
	var block ZmeyCoinBlock

	decoder := gob.NewDecoder(bytes.NewReader(serializedBlock))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (block ZmeyCoinBlock) NewGenesisBlock(coinbaseTransaction *Interface.Transaction) Interface2.Block {
	return NewBlock([]Interface.Transaction{*coinbaseTransaction}, []byte{}) //TODO: seems nonsense - check this
}

func (block ZmeyCoinBlock) ComputeHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.ComputeTransactionsHash(), timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = &hash[:]
}

func (block ZmeyCoinBlock) ComputeTransactionsHash() []byte {
	var transactionHashes [][]byte

	for _, tx := range block.Transactions {
		transactionHashes = append(transactionHashes, tx.ToBytes())
	}

	merkleTree := MerkleTree.NewMerkleTree(transactionHashes)

	return merkleTree.Root.Data
}

func (block ZmeyCoinBlock) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("ZmeyCoinBlock creation timestamp: %x\n",  time.Unix(block.Timestamp, 0)))
	buffer.WriteString(fmt.Sprintf("Hash: %x\n", block.Hash))
	buffer.WriteString(fmt.Sprintf("Prev. hash: %x\n", block.PrevBlockHash))
	buffer.WriteString(fmt.Sprintf("Data: %v\n", block.Transactions))

	return fmt.Sprintf("%v", buffer.String())
}

func (block ZmeyCoinBlock) Serialize() []byte {
	return Util.SerializeObject(block)
}
