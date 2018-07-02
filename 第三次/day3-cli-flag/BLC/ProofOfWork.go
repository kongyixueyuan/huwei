package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)
//256位hash里面前面至少有16个零
const targetBit = 16

type ProofOfWork struct {
	Block *Block //当前要验证的区块
	target *big.Int //大数据存储
}
func (proofOfWork *ProofOfWork) IsValid() bool{
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)
	if proofOfWork.target.Cmp(&hashInt)==1 {
		return true
	}
   return false
}
//数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int)  []byte{
	data := bytes.Join([][]byte{
		pow.Block.prevBlockHash,
		pow.Block.Data,
		IntToHex(pow.Block.Timestamp),
		IntToHex(int64(targetBit)),
		IntToHex(int64(nonce)),
		IntToHex(int64(pow.Block.Height)),
	},[]byte{})
	return  data
}
//行为
func (proofOfWork * ProofOfWork) Run()([]byte,int64)  {
	//1. 将Block的属性拼接成字节数组
	nonce := 0
	var hashInt big.Int  //存储新生成的HASH
	var hash [32]byte
	for {
		//准备数据
		dataBytes := proofOfWork.prepareData(nonce)
		hash = sha256.Sum256(dataBytes)
		fmt.Print("time:",time.Now() )
		fmt.Printf("\r%x\n :",hash)
        //将hash存储到hashInt

		hashInt.SetBytes(hash[:])
		//fmt.Println(hashInt)
		//判断hashInt是否小于Block里面的target //3. 判断hash有效性,如果满足条件，跳出循环
		if proofOfWork.target.Cmp(&hashInt) ==1 {
			break
		}
		nonce = nonce+1
	}



	return hash[:],int64(nonce)
}
//创建新的工作量证明
func NewProofOfWork(block *Block)  *ProofOfWork{
	//1. bit.Int对象
	//2.创建一个初始值为1的target
	target := big.NewInt(1)
	//3. 左移256 -target
	target=target.Lsh(target,256-targetBit)
	return &ProofOfWork{Block:block,target:target}

}
