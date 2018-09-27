package Interface

type BlockchainIterator interface {
	Next() Block
}