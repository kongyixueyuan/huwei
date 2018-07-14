package BLC

import (
	"math/big"
	"math"
	"bytes"
	"crypto/sha256"
	"fmt"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 16

type Hw_ProofOfWork struct {
	Hw_block  *Hw_Block
	Hw_target *big.Int
}

// 生成新的工作量证明
func Hw_NewProofOfWork(b *Hw_Block) *Hw_ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &Hw_ProofOfWork{b, target}
	return pow
}

// 准备挖矿hash数据
func (pow *Hw_ProofOfWork) Hw_PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Hw_block.Hw_PrevBlockHash,
		pow.Hw_block.Hw_HashTransactions(),
		IntToHex(pow.Hw_block.Hw_TimeStamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

// 执行工作量证明，返回nonce值和hash
func (pow *Hw_ProofOfWork) Hw_Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0
	for nonce < maxNonce {
		data := pow.Hw_PrepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		//if math.Remainder(float64(nonce),100000) == 0{
		//	fmt.Printf("\r%x",hash)
		//}
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.Hw_target) == -1 {
			break;
		} else {
			nonce++
		}
	}
	return nonce, hash[:]

}
