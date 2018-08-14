package MerkleTree

import (
		"crypto/sha256"
	"encoding/gob"
	"bytes"
	"log"
		)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	LeftLeaf  *MerkleNode
	RightLeaf *MerkleNode
	Data      []byte
}

func (merkleNode *MerkleNode) hashNode(data *[]byte) {
	if merkleNode.LeftLeaf == nil  && merkleNode.RightLeaf == nil {
		//The leaf holds just the double hash of the Transaction
		tempHashed := sha256.Sum256(*data)
		merkleNode.Data = tempHashed[:]
	} else {
		var temp []byte
		//If only one side is empty then the other Transaction is copied
		if merkleNode.LeftLeaf == nil{
			temp = append(merkleNode.RightLeaf.Data, merkleNode.RightLeaf.Data...)
		}
		if merkleNode.RightLeaf == nil{
			temp = append(merkleNode.LeftLeaf.Data, merkleNode.LeftLeaf.Data...)
		} else {
			temp = append(merkleNode.LeftLeaf.Data, merkleNode.RightLeaf.Data...)
		}
		tempHashed := sha256.Sum256(temp)
		merkleNode.Data = tempHashed[:]
	}
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode{
	merkleNode := MerkleNode{LeftLeaf: left, RightLeaf: right}
	merkleNode.hashNode(&data)
	return &merkleNode
}

func createTree (transactions [][]byte) *MerkleNode{
	if len(transactions) == 1 {
		return NewMerkleNode(nil, nil, transactions[0])
	}
	length := len(transactions)/2

	left := createTree(transactions[:length])
	right := createTree(transactions[length:])
	return NewMerkleNode(left, right, nil)
}
func NewMerkleTree(transactions [][]byte) *MerkleTree{
	length := len(transactions)
	if length %2 != 0 {
		transactions = append(transactions, transactions[length-1]) //In case the division isn't even, we'd like the right branch be smaller so we duplicate it.
	}
	return &MerkleTree{createTree(transactions)}
}

// DeserializeBlock deserialize a Block
func DeserializeMerkleNode(serializedMerkleNode []byte) *MerkleNode {
	var merkleNode MerkleNode

	decoder := gob.NewDecoder(bytes.NewReader(serializedMerkleNode))
	err := decoder.Decode(&merkleNode)
	if err != nil {
		log.Panic(err)
	}

	return &merkleNode
}
