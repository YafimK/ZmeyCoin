package Blockchain

import (
	"fmt"
	"errors"
	"encoding/hex"
	"bytes"
	"crypto/ecdsa"
	"log"
	"github.com/dgraph-io/badger"
	"ZmeyCoin/BlockChain/Interface"
	Interface2 "ZmeyCoin/Transaction/Interface"
	"ZmeyCoin/Transaction"
)
type ZmeyChain struct {
	//blocks []*ZmeyCoinBlock.ZmeyCoinBlock
	//transactions []*transaction.ZmeyTransaction //ZmeyTransaction pending to be "ZmeyCoinBlock'ed"
	//blocksCount int
	BlockTip     []byte
	BlockDb      *badger.DB
	ChainstateDb *badger.DB
}

func (blockchain *ZmeyChain) GetChainStateDb() *badger.DB {
	return blockchain.ChainstateDb
}

func (blockchain *ZmeyChain) AddBlock(transactions []*Interface2.Transaction) Interface.Block {
	var lastHash []byte

	err := blockchain.BlockDb.View(func(dbTransaction *badger.Txn) error {
		item ,err := dbTransaction.Get([]byte("l"))
		if err != nil {
			return errors.New(fmt.Sprintf("Error finding the ZmeyCoinBlock tip in the BlockDb: %v \n", err))
		}
		lastHash, err = item.Value()
		if err != nil {
			return errors.New(fmt.Sprintf("Error restoring the ZmeyCoinBlock tip from BlockDb: %v \n", err))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Errors during getting the ZmeyCoinBlock tip from the BlockDb: %v \n", err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = blockchain.BlockDb.Update(func(dbTransaction *badger.Txn) error {
		err := dbTransaction.Set(*newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis ZmeyCoinBlock into the BlockDb: %v \n", err))
		}
		err = dbTransaction.Set([]byte("l"), *newBlock.Hash)
		if err != nil {
			return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis ZmeyCoinBlock into the BlockDb: %v \n", err))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Errors during updating the ZmeyCoinBlock tip in the BlockDb: %v \n", err)
	}
	blockchain.BlockTip = *newBlock.Hash
	return newBlock
}

func (blockchain *ZmeyChain) MineBlock(transactions []*Interface2.Transaction) Interface.Block {
	var lastHash []byte

	for _, tx := range transactions {
		if blockchain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := blockchain.BlockDb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("l"))
		if err != nil {
			return err
		}
		lastHash, err = item.Value()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = blockchain.BlockDb.Update(func(txn *badger.Txn) error {
		err := txn.Set(*newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = txn.Set([]byte("l"), *newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		blockchain.BlockTip = *newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

func (blockchain *ZmeyChain) initBlockDb() (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/zmeyCoin/blocks"
	opts.ValueDir = "/tmp/zmeyCoin/blocks"
	blockDb, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return blockDb, err}

func (blockchain *ZmeyChain) initChainstateDb() (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/zmeyCoin/Chainstate"
	opts.ValueDir = "/tmp/zmeyCoin/Chainstate"
	chainstateDb, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return  chainstateDb, err}

func (blockchain *ZmeyChain) Iterator() Interface.BlockchainIterator {
	return NewBlockchainIterator(blockchain)
}

func (blockchain *ZmeyChain) SignTransaction(tx *Interface2.Transaction, privateKey ecdsa.PrivateKey) {
	previousTransactions := make(map[string]Transaction.ZmeyTransaction)

	for _, transactionInput := range tx.GetTransactionInputs() {
		previousTransaction, err := blockchain.FindTransaction(transactionInput.PrevTransactionHash)
		if err != nil {
			log.Panicf("error looking for ")
		}
		previousTransactions[hex.EncodeToString(previousTransaction.GetHash())] = previousTransaction
	}

	tx.(Transaction.ZmeyTransaction).SignTransaction(privateKey, previousTransactions)}

func (blockchain *ZmeyChain) FindUnspentTransactionOutputs() map[string]Interface2.TransactionOutputs {
	UTXO := make(map[string]Transaction.TxOutputs)
	spentTXOs := make(map[string][]int)
	bci := blockchain.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.GetHash())

		Outputs:
			for outIdx, out := range tx.TransactionOutputs {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs  = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.TransactionInputs {
					inTxID := hex.EncodeToString(in.PrevTransactionHash)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.PrevTxOutIndex)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
	}

//func NewBlockChain() *ZmeyChain {
//	blockDb, err := initBlockDb()
//	chainstateDb, err := initChainstateDb()
//	//TODO:  return the defer of closing the db back in game
//
//	var tip []byte
//
//	err = blockDb.Update(func(tx *badger.Txn) error {
//		item, err := tx.Get([]byte("l"))
//		if err == badger.ErrKeyNotFound {
//			genesis := NewGenesisBlock(Transaction.NewCoinbaseTransaction())
//			err = tx.Set(*genesis.Hash, genesis.Serialize())
//			if err != nil {
//				return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis ZmeyCoinBlock into the BlockDb: %v \n", err))
//			}
//			err = tx.Set([]byte("l"), *genesis.Hash)
//			if err != nil {
//				return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis ZmeyCoinBlock into the BlockDb: %v \n", err))
//			}
//			tip = *genesis.Hash
//		} else if err != nil {
//			return errors.New(fmt.Sprintf("We had some issues finding the ZmeyCoinBlock tip in the BlockDb: %v \n", err))
//		} else {
//			tip, err = item.Value()
//			if err != nil {
//				return errors.New(fmt.Sprintf("We had some issues restoring the ZmeyCoinBlock tip from BlockDb: %v \n", err))
//			}
//		}
//
//		return nil
//	})
//	if err != nil {
//		log.Fatalf("We had some issues finding the ZmeyCoinBlock tip in the BlockDb: %v \n", err)
//	}
//	return &ZmeyChain{BlockTip: tip, BlockDb: blockDb, ChainstateDb: chainstateDb}
//}

func (blockchain *ZmeyChain) PrintBlockChain() {
	fmt.Println("*** ZmeyChain ***")
	//for index, curBlock := range blockchain.blocks {
	//	fmt.Printf("%v ZmeyCoinBlock\n",index)
	//	fmt.Println(curBlock)
	//}
}

func (blockchain *ZmeyChain) AddTransaction() {

}

func (blockchain *ZmeyChain) FindTransaction(targetTransactionHash []byte) (Interface2.Transaction, error){
	blockchainIterator := blockchain.Iterator()
	
	for currentBlock := blockchainIterator.Next(); currentBlock != nil; {
		for _, tx := range currentBlock.Transactions {
			if bytes.Compare(tx.GetHash(), targetTransactionHash) == 0 {
				return *tx, nil
			}
		}
	}

	return Transaction.ZmeyTransaction{}, errors.New("transaction is not found")
}

func (blockchain *ZmeyChain) VerifyTransaction(targetTransaction *Interface2.Transaction ) bool {
	prevTXs := make(map[string]Transaction.ZmeyTransaction)

	for _, vin := range targetTransaction.TransactionInputs {
		prevTX, err := blockchain.FindTransaction(vin.PrevTransactionHash)
		if err != nil {
			log.Fatalf("failed finding suitable transaction: %v\n", err)
		}
		prevTXs[hex.EncodeToString(prevTX.GetHash())] = prevTX
	}

	return targetTransaction.Verify(prevTXs)
}

