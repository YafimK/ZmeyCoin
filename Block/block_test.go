package Block

import (
	"testing"
	"bytes"
)

func TestBlockCreation(t *testing.T){
	blockData := "test Block data"
	result := New(blockData, []byte{})
	if result == nil {
		t.Errorf("Generated nil block1")
	}
	if !bytes.Equal(result.Data, []byte(blockData)) {
		t.Errorf("The Block data doesn't match the wanted data")
	}
}
