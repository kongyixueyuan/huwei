package BLC

/**
    定义一个区块的结构
 */

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
)

type Block struct {
	//1.区块高度
	Height int64
	//2.上一个区块的HASH
	prevBlockHash []byte
	//3.交易数据
	Data []byte
	//4.时间戳
	Timestamp int64
	//5.hash
	Hash []byte
	//6. nonce
	Nonce int64
}

//1.创建新的区块 *代表指针  @引用变量的地址
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{Height: height, prevBlockHash: prevBlockHash, Data: []byte(data), Timestamp: time.Now().Unix(), Hash: nil,Nonce:0}
	//根据区块数据生成当前区块的Hash
	//调用工作量证明方法并且返回有效的hash和nonce
	pow := NewProofOfWork(block)
	//挖矿验证
	hash,nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//2 单独写一个方法，生成创世区块
func CreateGenesisBlock(data string) *Block{
	return NewBlock(data,1,[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

}

//2.将区块内容设置成Hash  base 后面是进制进行转换 2-36
func (block *Block) setHash() {
	// 1 .Heigh [] byte
	heightBytes := IntToHex(block.Height)
	//fmt.Println("heightBytes ", heightBytes)
	// 2. 时间戳转成 []byte 2~36
	timeString := strconv.FormatInt(block.Timestamp, 2)
    timeBytes := []byte(timeString)
	//fmt.Println("timeBytes ", timeBytes)
	//3.拼接所有属性
	blockBytes := bytes.Join([][]byte{heightBytes,block.prevBlockHash,block.Data,timeBytes,block.Hash},[]byte{})
	//生成Hash
	hash := sha256.Sum256(blockBytes)
	block.Hash = hash[:]

}

