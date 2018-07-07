package BLC

/**
    定义一个区块的结构
 */

import (
	"time"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"fmt"
)

type Block struct {
	//1.区块高度
	Height int64
	//2.上一个区块的HASH
	PrevBlockHash []byte
	//3.交易数据
	Txs []*Transaction
	//4.时间戳
	Timestamp int64
	//5.hash
	Hash []byte
	//6. nonce
	Nonce int64
}

//需要将Txs转换成byte
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

//序列化，把区块对象转成[]byte
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化，将字节数组转成对象
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//1.创建新的区块 *代表指针  @引用变量的地址
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{Height: height,PrevBlockHash: prevBlockHash, Txs: txs, Timestamp: time.Now().Unix(), Hash: nil, Nonce: 0}
	//调用工作量证明方法并且返回有效的hash和nonce
	pow := NewProofOfWork(block)
	//挖矿验证
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	fmt.Println()
	return block
}

//2 单独写一个方法，生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

}
