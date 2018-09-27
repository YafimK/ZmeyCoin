package Interface

type TransactionInput interface {
	UsesKey (pubKeyHash []byte) bool
	String() string
	GetPrevTransactionHash() []byte
	GetPrevTxOutIndex() int
}
