package MerkleTree



type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	LeftLeaf *MerkleNode
	RightLeaf *MerkleNode
	data []byte
}