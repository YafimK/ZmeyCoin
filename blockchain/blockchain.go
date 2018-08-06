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
const dbFile = "blockchain.dat"

type Blockchain struct {
	//blocks []*Block.Block
	//transactions []*transaction.Transaction //Transaction pending to be "Block'ed"
	//blocksCount int
	BlockTip *[]byte
	db *badger.DB

}

func (blockchain *Blockchain) AddBlock(transactions []*transaction.Transaction) {
	prevBlock := blockchain.blocks[blockchain.blocksCount - 1]
	newBlock := Block.New(transactions, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
	blockchain.blocksCount++
}

func (blockchain *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	//TODO: gather all possible transactions and create a new Block
	blockchain.AddBlock(transactions)
}

//we need to init the blockchain with genesis Block
func New() *Blockchain {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/zmeyCoin"
	opts.ValueDir = "/tmp/zmeyCoin"
	var tip []byte
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *badger.Txn) error {
		item, err := tx.Get([]byte("l"))
		if  err == badger.ErrKeyNotFound {
			genesis := NewGenesisBlock()
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
			//if err != nil {
			//	log.Println("We had some issues restoring the Block tip from db", err)
			//}
		}

		return err
	})
	log.Fatalf("We had some issues finding the Block tip in the db: %v \n", err)

	return &Blockchain{&tip, db}
}
func NewGenesisBlock() *Block.Block {

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
	bci := blockchain.Iterator()
	
	for currentBlock := bci.Next(); currentBlock != nil; {

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

func (blockchain *Blockchain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	previousTransactions := make(map[string]transaction.Transaction)

	for _, transactionInput := range tx.TransactionInputs {
		previousTransaction, err := blockchain.FindTransaction(transactionInput.PrevTransactionHash)
		if err != nil {
			log.Panicf("error looking for ")
		}
		previousTransactions[hex.EncodeToString(previousTransaction.GetHash())] = previousTransaction
	}

	tx.SignTransaction(privKey, previousTransactions)
}

