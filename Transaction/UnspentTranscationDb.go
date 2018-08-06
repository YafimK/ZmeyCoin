package Transaction

import (
	"ZmeyCoin/BlockChain"
	"encoding/hex"
	"log"
	"github.com/dgraph-io/badger"
		"ZmeyCoin/Block"
	)

type UnspentTransactionIndex struct {
	Blockchain *BlockChain.Blockchain

}

func (u *UnspentTransactionIndex) FindUnspentOutputsPerPubKeyHash(pubKeyHash []byte) []TransactionOutput {
	var UnspentOutputs []TransactionOutput
	db := u.Blockchain.ChainstateDb
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}
			outs := DeserializeOutputs(v)
			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UnspentOutputs = append(UnspentOutputs, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return UnspentOutputs
}


func (u *UnspentTransactionIndex) CountTransactions() int {
	db := u.Blockchain.ChainstateDb
	counter := 0

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			counter++
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return counter
}

func (u *UnspentTransactionIndex) Reindex() {
	db := u.Blockchain.ChainstateDb
	UTXO := u.Blockchain.FindUnspentTransactionOutputs()
	err := db.Update(func(txn *badger.Txn) error {
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = txn.Set(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (u *UnspentTransactionIndex) FindSpendableOutputs(publicKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.ChainstateDb
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)
			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(publicKeyHash) && accumulated < amount {
					accumulated += out.value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

func (u *UnspentTransactionIndex) Update(block *Block.Block) {
	db := u.Blockchain.ChainstateDb
	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _, vin := range tx.TransactionInputs {
					updatedOuts := TransactionOutputs{}
					item, err := txn.Get(vin.PrevTransactionHash)
					if err != nil {return err}
					outsBytes,err := item.Value()
					outs := DeserializeOutputs(outsBytes)
					for outIdx, out := range outs.Outputs {
						if outIdx != vin.PrevTxOutIndex {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}
					if len(updatedOuts.Outputs) == 0 {
						err := txn.Delete(vin.PrevTransactionHash)
						if err != nil {
							log.Panic(err)
						}
					} else {
						err := txn.Set(vin.PrevTransactionHash, updatedOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}
			newOutputs := TransactionOutputs{}
			for _, out := range tx.TransactionOutputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			err := txn.Set(tx.GetHash(), newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}