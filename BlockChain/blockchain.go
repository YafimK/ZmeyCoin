package BlockChain

import (
	"ZmeyCoin/Block"
	"fmt"
	"ZmeyCoin/Transaction"
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
	blockDb  *badger.DB
	chainstateDb  *badger.DB

}

func (blockchain *Blockchain) AddBlock(transactions []*Transaction.Transaction) *Block.Block{
	var lastHash []byte

	err := blockchain.blockDb.View(func(dbTransaction *badger.Txn) error {
		item ,err := dbTransaction.Get([]byte("l"))
		if err != nil {
			return errors.New(fmt.Sprintf("Error finding the Block tip in the blockDb: %v \n", err))
		}
		lastHash, err = item.Value()
		if err != nil {
			return errors.New(fmt.Sprintf("Error restoring the Block tip from blockDb: %v \n", err))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Errors during getting the Block tip from the blockDb: %v \n", err)
	}

	newBlock := Block.NewBlock(transactions, lastHash)

	err = blockchain.blockDb.Update(func(dbTransaction *badger.Txn) error {
		err := dbTransaction.Set(*newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the blockDb: %v \n", err))
		}
		err = dbTransaction.Set([]byte("l"), *newBlock.Hash)
		if err != nil {
			return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the blockDb: %v \n", err))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Errors during updating the Block tip in the blockDb: %v \n", err)
	}
	blockchain.BlockTip = *newBlock.Hash
	return newBlock
}

func (blockchain *Blockchain) MineBlock(transactions []*Transaction.Transaction) *Block.Block {
	var lastHash []byte

	for _, tx := range transactions {
		if blockchain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := blockchain.blockDb.View(func(txn *badger.Txn) error {
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

	err = blockchain.blockDb.Update(func(txn *badger.Txn) error {
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
	blockDb, err := initBlockDb()
	chainstateDb, err := initChainstateDb()
	//TODO:  return the defer of closing the db back in game

	var tip []byte

	err = blockDb.Update(func(tx *badger.Txn) error {
		item, err := tx.Get([]byte("l"))
		if err == badger.ErrKeyNotFound {
			genesis := NewGenesisBlock(Transaction.NewCoinbaseTransaction())
			err = tx.Set(*genesis.Hash, genesis.Serialize())
			if err != nil {
				return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the blockDb: %v \n", err))
			}
			err = tx.Set([]byte("l"), *genesis.Hash)
			if err != nil {
				return errors.New(fmt.Sprintf("We had some issues inserting hash of genesis Block into the blockDb: %v \n", err))
			}
			tip = *genesis.Hash
		} else if err != nil {
			return errors.New(fmt.Sprintf("We had some issues finding the Block tip in the blockDb: %v \n", err))
		} else {
			tip, err = item.Value()
			if err != nil {
				return errors.New(fmt.Sprintf("We had some issues restoring the Block tip from blockDb: %v \n", err))
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("We had some issues finding the Block tip in the blockDb: %v \n", err)
	}
	return &Blockchain{BlockTip: tip, blockDb: blockDb, chainstateDb: chainstateDb}
}

func initBlockDb() (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/zmeyCoin/blocks"
	opts.ValueDir = "/tmp/zmeyCoin/blocks"
	blockDb, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return blockDb, err
}

func initChainstateDb() (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/zmeyCoin/Chainstate"
	opts.ValueDir = "/tmp/zmeyCoin/Chainstate"
	chainstateDb, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return  chainstateDb, err
}

func NewGenesisBlock(coinbaseTransaction *Transaction.Transaction) *Block.Block {

	return Block.NewBlock([]*Transaction.Transaction{coinbaseTransaction}, []byte{})
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

func (blockchain *Blockchain) FindTransaction(targetTransactionHash []byte) (Transaction.Transaction, error){
	blockchainIterator := blockchain.Iterator()
	
	for currentBlock := blockchainIterator.Next(); currentBlock != nil; {
		for _, tx := range currentBlock.Transactions {
			if bytes.Compare(tx.GetHash(), targetTransactionHash) == 0 {
				return *tx, nil
			}
		}
	}

	return Transaction.Transaction{}, errors.New("transaction is not found")
}

func (blockchain *Blockchain) VerifyTransaction(targetTransaction *Transaction.Transaction) bool {
	prevTXs := make(map[string]Transaction.Transaction)

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

func (blockchain *Blockchain) SignTransaction(tx *Transaction.Transaction, privateKey ecdsa.PrivateKey) {
	previousTransactions := make(map[string]Transaction.Transaction)

	for _, transactionInput := range tx.TransactionInputs {
		previousTransaction, err := blockchain.FindTransaction(transactionInput.PrevTransactionHash)
		if err != nil {
			log.Panicf("error looking for ")
		}
		previousTransactions[hex.EncodeToString(previousTransaction.GetHash())] = previousTransaction
	}

	tx.SignTransaction(privateKey, previousTransactions)
}

func (blockchain *Blockchain) FindUnspentTransactionOutputs()  map[string]Transaction.TransactionOutputs {
	UTXO := make(map[string]Transaction.TransactionOutputs)
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