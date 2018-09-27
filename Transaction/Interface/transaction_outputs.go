package Interface

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TransactionOutputs collects TXOutput
type TransactionOutputs interface {
	SerializeOutputs() []byte
	GetOutputs() []TransactionOutput
}

func DeserializeOutputs(data []byte) TransactionOutputs {
	var outputs TransactionOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
