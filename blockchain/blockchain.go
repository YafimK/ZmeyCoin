package blockchain

import (
	"ZmeyCoin/Block"
	"fmt"
	"ZmeyCoin/transaction"
		"errors"
	"encoding/hex"
	"bytes"
	"crypto/ecdsa"
	"log"
	"github.com/dgraph-io/badger"
	)
type Blockchain struct {
	//blocks []*Block.Block
	//transactions []*transaction.Transaction //Transaction pending to be "Block'ed"
	//blocksCount int
	BlockTip []byte
	db *badger.DB

}

func (blockchain *Blockchain) AddBlock(transactions []*transaction.Transaction) *Block.Block{
	var lastHash []byte

	err := blockchain.db.View(func(dbTransaction *badger.Txn) error {
		item ,err := dbTransaction.Get([]byte("l"))
		if err != nil {
			return errors.New(fmt.Sprintf("Error finding the Block tip in the db: %v \n", err))
		}
		lastHash, err = item.Value()
		if err != nil {
			return errors.New(fmt.Sprintf("Error restoring the Block tip from db: %v \n", err))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Errors during getting the Block tip from the db: %v \n", err)
	}

	newBlock := Block.NewBlock(transactions, lastHash)

	err = blockchain.db.Update(func(dbTransaction *badger.Txn) error {
		err := dbTransaction.Set(*newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the db: %v \n", err))
		}
		err = dbTransaction.Set([]byte("l"), *newBlock.Hash)
		if err != nil {
			return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the db: %v \n", err))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Errors during updating the Block tip in the db: %v \n", err)
	}
	blockchain.BlockTip = *newBlock.Hash
	return newBlock
}

func (blockchain *Blockchain) MineBlock(transactions []*transaction.Transaction) *Block.Block {
	var lastHash []byte

	for _, tx := range transactions {
		if blockchain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := blockchain.db.View(func(txn *badger.Txn) error {
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

	newBlock := Block.NewBlock(transactions, lastHash)

	err = blockchain.db.Update(func(txn *badger.Txn) error {
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

//we need to init the blockchain with genesis Block
func New() *Blockchain {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/zmeyCoin/blocks"
	opts.ValueDir = "/tmp/zmeyCoin/blocks"
	var tip []byte
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *badger.Txn) error {
		item, err := tx.Get([]byte("l"))
		if  err == badger.ErrKeyNotFound {
			genesis := NewGenesisBlock(transaction.NewCoinbaseTransaction())
			err = tx.Set(*genesis.Hash, genesis.Serialize())
			if err != nil {
				return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the db: %v \n", err))
			}
			err = tx.Set([]byte("l"), *genesis.Hash)
			if err != nil {
				return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the db: %v \n", err))
			}
			tip = *genesis.Hash
		} else if err != nil {
			return errors.New(fmt.Sprintf("We had some issues finding the Block tip in the db: %v \n", err))
		} else {
			tip, err = item.Value()
			if err != nil {
				return errors.New(fmt.Sprintf("We had some issues restoring the Block tip from db: %v \n", err))
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("We had some issues finding the Block tip in the db: %v \n", err)
	}

	return &Blockchain{tip, db}
}
func NewGenesisBlock(coinbaseTransaction *transaction.Transaction) *Block.Block {

	return Block.NewBlock([]*transaction.Transaction{coinbaseTransaction}, []byte{})
}

func (blockchain *Blockchain) PrintBlockChain() {
	fmt.Println("*** Blockchain ***")
	//for index, curBlock := range blockchain.blocks {
	//	fmt.Printf("%v Block\n",index)
	//	fmt.Println(curBlock)
	//}
}

func (blockchain *Blockchain) AddTransaction() {

}

func (blockchain *Blockchain) FindTransaction(targetTransactionHash []byte) (transaction.Transaction, error){
	blockchainIterator := blockchain.Iterator()
	
	for currentBlock := blockchainIterator.Next(); currentBlock != nil; {
		for _, tx := range currentBlock.Transactions {
			if bytes.Compare(tx.GetHash(), targetTransactionHash) == 0 {
				return *tx, nil
			}
		}
	}

	return transaction.Transaction{}, errors.New("transaction is not found")
}

func (blockchain *Blockchain) VerifyTransaction(targetTransaction *transaction.Transaction) bool {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range targetTransaction.TransactionInputs {
		prevTX, err := blockchain.FindTransaction(vin.PrevTransactionHash)
		if err != nil {
			log.Fatalf("failed finding suitable transaction: %v\n", err)
		}
		prevTXs[hex.EncodeToString(prevTX.GetHash())] = prevTX
	}

	return targetTransaction.Verify(prevTXs)
}

func (blockchain *Blockchain) Iterator()  *BlockchainIterator{
	return NewBlockchainIterator(blockchain)
}

func (blockchain *Blockchain) SignTransaction(tx *transaction.Transaction, privateKey ecdsa.PrivateKey) {
	previousTransactions := make(map[string]transaction.Transaction)

	for _, transactionInput := range tx.TransactionInputs {
		previousTransaction, err := blockchain.FindTransaction(transactionInput.PrevTransactionHash)
		if err != nil {
			log.Panicf("error looking for ")
		}
		previousTransactions[hex.EncodeToString(previousTransaction.GetHash())] = previousTransaction
	}

	tx.SignTransaction(privateKey, previousTransactions)
}

