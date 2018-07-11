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

type HW_ProofOfWork struct {
	hw_block  *HW_Block
	hw_target *big.Int
}

// 生成新的工作量证明
func HW_NewProofOfWork(b *HW_Block) *HW_ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &HW_ProofOfWork{b, target}
	return pow
}

// 准备挖矿hash数据
func (pow *HW_ProofOfWork) HW_PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.hw_block.Hw_PrevBlockHash,
		pow.hw_block.HW_HashTransactions(),
		IntToHex(pow.hw_block.Hw_TimeStamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

// 执行工作量证明，返回nonce值和hash
func (pow *HW_ProofOfWork) HW_Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0
	for nonce < maxNonce {
		data := pow.HW_PrepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		//if math.Remainder(float64(nonce),100000) == 0{
		//	fmt.Printf("\r%x",hash)
		//}
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.hw_target) == -1 {
			break;
		} else {
			nonce++
		}
	}
	return nonce, hash[:]

}
