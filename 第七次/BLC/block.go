package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)

type Hw_Block struct {
	Hw_TimeStamp     int64
	Hw_Transactions   []*Hw_Transaction
	Hw_PrevBlockHash []byte
	Hw_Hash          []byte
	Hw_Nonce         int
	Hw_Height        int
}
// 生成新的区块
func Hw_NewBlock(transactions []*Hw_Transaction, prevBlockHash []byte, height int) *Hw_Block {
	// 生成新的区块对象
	block := &Hw_Block{
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
		height,
	}
	// 挖矿

	pow := Hw_NewProofOfWork(block)
	nonce,hash :=pow.Hw_Run()

	block.Hw_Nonce = nonce
	block.Hw_Hash = hash[:]

	return block

}

// 将交易进行hash
func (b Hw_Block) Hw_HashTransactions() []byte {
	var transactions [][]byte
	// 获取交易真实内容
	for _,tx := range b.Hw_Transactions{
		transactions = append(transactions,tx.Hw_Serialize())
	}
	//txHash := sha256.Sum256(bytes.Join(transactions,[]byte{}))
	mTree := Hw_NewMerkelTree(transactions)
	return mTree.Hw_RootNode.Hw_Data
}
// 新建创世区块
func Hw_NewGenesisBlock(coinbase *Hw_Transaction) *Hw_Block  {
	return Hw_NewBlock([]*Hw_Transaction{coinbase},[]byte{},1)
}

// 序列化
func (b *Hw_Block) Hw_Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func Hw_DeserializeBlock(d []byte) *Hw_Block {
	var block Hw_Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
// 打印区块内容
func (block Hw_Block) String()  {
	fmt.Println("\n==============")
	fmt.Printf("Height:\t%d\n", block.Hw_Height)
	fmt.Printf("PrevBlockHash:\t%x\n", block.Hw_PrevBlockHash)
	fmt.Printf("Timestamp:\t%s\n", time.Unix(block.Hw_TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	fmt.Printf("Hash:\t%x\n", block.Hw_Hash)
	fmt.Printf("Nonce:\t%d\n", block.Hw_Nonce)
	fmt.Println("Txs:")

	for _, tx := range block.Hw_Transactions {
		tx.String()
	}
	fmt.Println("==============")
}
