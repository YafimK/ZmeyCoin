package Transaction

import (
	"ZmeyCoin/Transaction/Interface"
	"ZmeyCoin/Util"
)

// TxOutputs is a collection of TXOutput
type TxOutputs struct {
	Outputs []Interface.TransactionOutput
}

func (outs TxOutputs) GetOutputs() []Interface.TransactionOutput {
	return outs.Outputs
}

func (outs TxOutputs) SerializeOutputs() []byte {
	return Util.SerializeObject(outs)
}
