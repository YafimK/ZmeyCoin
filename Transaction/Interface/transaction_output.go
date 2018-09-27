package Interface

//TransactionOutput
type TransactionOutput interface {
	lockOutputByAddress(address []byte)
	IsLockedWithKey(pubKeyHash []byte) bool
	GetValue() interface{}
}
