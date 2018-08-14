package MerkleTree
import (
	"encoding/hex"
	"fmt"
	"testing"
	"strings"
)

func TestNewMerkleNode(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}

	// Level 1

	n1 := NewMerkleNode(nil, nil, data[0])
	n2 := NewMerkleNode(nil, nil, data[1])
	n3 := NewMerkleNode(nil, nil, data[2])
	n4 := NewMerkleNode(nil, nil, data[2])

	// Level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// Level 3
	n7 := NewMerkleNode(n5, n6, nil)


	if strings.Compare("64b04b718d8b7c5b6fd17f7ec221945c034cfce3be4118da33244966150c4bd4", hex.EncodeToString(n5.Data)) != 0 {
		t.Errorf("Level 1 hash 1 is not correct - expected got %v\n", hex.EncodeToString(n5.Data))
	}
	if strings.Compare("08bd0d1426f87a78bfc2f0b13eccdf6f5b58dac6b37a7b9441c1a2fab415d76c", hex.EncodeToString(n6.Data)) != 0 {
		t.Errorf("Level 1 hash 2 is not correct - expected got %v\n", hex.EncodeToString(n6.Data))
	}
	if strings.Compare("4e3e44e55926330ab6c31892f980f8bfd1a6e910ff1ebc3f778211377f35227e", hex.EncodeToString(n7.Data)) != 0 {
		t.Errorf("Root hash is not correct - expected got %v\n", hex.EncodeToString(n7.Data))
	}
}

func TestNewMerkleTree(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}
	// Level 1
	n1 := NewMerkleNode(nil, nil, data[0])
	n2 := NewMerkleNode(nil, nil, data[1])
	n3 := NewMerkleNode(nil, nil, data[2])
	n4 := NewMerkleNode(nil, nil, data[2])

	// Level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// Level 3
	n7 := NewMerkleNode(n5, n6, nil)

	rootHash := fmt.Sprintf("%x", n7.Data)
	mTree := NewMerkleTree(data)

	if strings.Compare(rootHash, fmt.Sprintf("%x", mTree.Root.Data)) != 0 {
		t.Errorf("Merkle tree root hash is not correct - expected \n%v\n got\n %v\n", rootHash, fmt.Sprintf("%x", mTree.Root.Data))
	}
}