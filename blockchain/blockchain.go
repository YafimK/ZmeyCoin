package blockchain

import (
	"ZmeyCoin/block"
	"fmt"
	"ZmeyCoin/transaction"
		"errors"
	"encoding/hex"
	"bytes"
	"crypto/ecdsa"
	"log"
)

type Blockchain struct {
	blocks []*block.Block
	transactions []*transaction.Transaction //Transaction pending to be "block'ed"
	blocksCount int
}

func (blockchain *Blockchain) AddBlock(transactions []*transaction.Transaction) {
	prevBlock := blockchain.blocks[blockchain.blocksCount - 1]
	newBlock := block.New(transactions, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
	blockchain.blocksCount++
}

func (blockchain *Blockchain) MineBlock() {
	//TODO: gather all possible transactions and create a new block

	blockchain.AddBlock([]*transaction.Transaction{})
}

//we need to init the blockchain with genesis block
func New() *Blockchain {
	newBlockchain := Blockchain{}
	newBlockchain.blocks = append(newBlockchain.blocks,
		block.New([]*transaction.Transaction{transaction.NewCoinbaseTransaction()}, []byte{}))
	newBlockchain.blocksCount++
	return &newBlockchain
}

func (blockchain *Blockchain) PrintBlockChain() {
	fmt.Println("*** Blockchain ***")
	for index, curBlock := range blockchain.blocks {
		fmt.Printf("%v block\n",index)
		fmt.Println(curBlock)
	}
}

func (blockchain *Blockchain) AddTransaction() {

}

func (blockchain *Blockchain) FindTransaction(targetTransactionHash []byte) (transaction.Transaction, error){
	bci := blockchain.Iterator()
	
	for {
		currentBlock := bci.Next()
	
		for _, tx := range currentBlock.Transactions {
			if bytes.Compare(tx.GetHash(), targetTransactionHash) == 0 {
				return *tx, nil
			}
		}
	
		if len(currentBlock.PrevBlockHash) == 0 {
			break
		}
	}

	return transaction.Transaction{}, errors.New("transaction is not found")
}

func (blockchain *Blockchain) VerifyTransaction(targetTransaction *transaction.Transaction) bool {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range targetTransaction.TransactionInputs {
		prevTX, err := blockchain.FindTransaction(vin.PrevTransactionHash)
		if err != nil {
			
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