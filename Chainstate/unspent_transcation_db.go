package Chainstate

import (
	"ZmeyCoin/BlockChain/Interface"
	"ZmeyCoin/Transaction"
	Interface2 "ZmeyCoin/Transaction/Interface"
	"encoding/hex"
	"github.com/dgraph-io/badger"
	"log"
)

type ZmeyCoinChainstate struct {
	Blockchain Interface.Blockchain
}

func (u *ZmeyCoinChainstate) FindUnspentOutputsPerPubKeyHash(pubKeyHash []byte) []Interface2.TransactionOutput {
	var UnspentOutputs []Interface2.TransactionOutput
	db := u.Blockchain.GetChainStateDb()
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
			outs := Interface2.DeserializeOutputs(v)
			for _, out := range outs.GetOutputs() {
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

func (u *ZmeyCoinChainstate) CountTransactions() int {
	db := u.Blockchain.GetChainStateDb()
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

func (u *ZmeyCoinChainstate) Reindex() {
	db := u.Blockchain.GetChainStateDb()
	UTXO := u.Blockchain.FindUnspentTransactionOutputs()
	err := db.Update(func(txn *badger.Txn) error {
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = txn.Set(key, outs.SerializeOutputs())
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

func (u *ZmeyCoinChainstate) FindSpendableOutputs(publicKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.GetChainStateDb()
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
			outs := Interface2.DeserializeOutputs(v)
			for outIdx, out := range outs.GetOutputs() {
				if out.IsLockedWithKey(publicKeyHash) && accumulated < amount {
					accumulated += out.GetValue().(int)
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

func (u *ZmeyCoinChainstate) Update(block Interface.Block) {
	db := u.Blockchain.GetChainStateDb()
	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range block.GetTransactions() {
			if tx.IsCoinbase() == false {
				for _, vin := range tx.GetTransactionInputs() {
					updatedOuts := Transaction.TxOutputs{}
					item, err := txn.Get(vin.PrevTransactionHash)
					if err != nil {return err}
					outsBytes,err := item.Value()
					outs := Transaction.DeserializeOutputs(outsBytes)
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
			newOutputs := tr.TxOutputs{}
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