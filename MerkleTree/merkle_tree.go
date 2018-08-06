package MerkleTree

import (
		sha256 "crypto/sha256"
	"encoding/gob"
	"bytes"
	"log"
	"ZmeyCoin/util"
	)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	LeftLeaf  *MerkleNode
	RightLeaf *MerkleNode
	Data      *[]byte
}

func (merkleNode *MerkleNode) hashNode(data *[]byte) {
	if merkleNode.LeftLeaf, merkleNode.RightLeaf == nil {
		//The leaf holds just the double hash of the transaction
		tempHashed := sha256.Sum256(*data)[:]
		merkleNode.Data = &tempHashed
	} else {
		var temp []byte
		//If only one side is empty then the other transaction is copied
		if merkleNode.LeftLeaf == nil{
			temp = append(util.SerializeObject(merkleNode.RightLeaf), util.SerializeObject(merkleNode.RightLeaf)...)
		}
		if merkleNode.RightLeaf == nil{
			temp = append(util.SerializeObject(merkleNode.LeftLeaf), util.SerializeObject(merkleNode.LeftLeaf)...)
		} else {
			temp = append(util.SerializeObject(merkleNode.LeftLeaf), util.SerializeObject(merkleNode.RightLeaf)...)
		}
		tempHashed := sha256.Sum256(temp)[:]
		merkleNode.Data = &tempHashed
	}
}

func NewMerkleNode(left, right *MerkleNode, data *[]byte) *MerkleNode{
	merkleNode := MerkleNode{LeftLeaf: left, RightLeaf: right}
	merkleNode.hashNode(data)
	return &merkleNode
}

func createTree (transactions *[][]byte) *MerkleNode{
	if len(*transactions) == 1 {
		return NewMerkleNode(nil, nil, &(*transactions)[0])
	} else if len(*transactions) == 0 {
		return nil
	}
	length := len(*transactions)/2
	left := createTree(&(*transactions)[:length])
	right := createTree(&(*transactions)[length:])
	return &MerkleNode{LeftLeaf:left,
	RightLeaf: right, Data:nil}
}
func NewMerkleTree(transactions *[][]byte) *MerkleTree{
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
