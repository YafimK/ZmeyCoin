package block

import (
	"testing"
	"bytes"
)

func TestBlockCreation(t *testing.T){
	blockData := "test block data"
	result := New(blockData, []byte{})
	if result == nil {
		t.Errorf("Generated nil block1")
	}
	if !bytes.Equal(result.Data, []byte(blockData)) {
		t.Errorf("The block data doesn't match the wanted data")
	}
}
