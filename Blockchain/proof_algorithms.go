package Blockchain

import (
	"math/big"
	"crypto/sha256"
	"bytes"
	"math"
	"ZmeyCoin/Util"
)

var ProofOfWorkTargetBits = 16 //TODO: change this with time
const MaxNonceValue = math.MaxInt64

type ProofOfWork struct {
	BlockTip         *ZmeyCoinBlock
	TargetDifficulty *big.Int
}

func NewProofOfWork(block *ZmeyCoinBlock) *ProofOfWork{
	targetDifficulty := big.NewInt(1)
	targetDifficulty.Lsh(targetDifficulty, uint(256 - ProofOfWorkTargetBits))

	return &ProofOfWork{BlockTip: block, TargetDifficulty: targetDifficulty}
	}


func (proof *ProofOfWork) CalculateProof() (int, *[]byte){
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for ;nonce < MaxNonceValue; nonce++ {
		data := proof.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(proof.TargetDifficulty) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, &hash[:]
}

//Concat all the needed that data for hashing
func (proof *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			proof.BlockTip.PrevBlockHash,
			proof.BlockTip.ComputeTransactionsHash(),
			Util.IntToHex(proof.BlockTip.Timestamp),
			Util.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

}

func (proof *ProofOfWork) Verify() bool {
	var hashInt big.Int

	data := proof.prepareData(proof.BlockTip.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(proof.TargetDifficulty) == -1

	return isValid
}