package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)
/**
   区块对象
 */
type HW_Block struct {
	Hw_TimeStamp     int64
	Hw_Transactions  []*HW_Transaction
	Hw_PrevBlockHash []byte
	Hw_Hash          []byte
	Hw_Nonce         int
	Hw_Height        int
}
// 生成新的区块
func HW_NewBlock(transactions []*HW_Transaction, prevBlockHash []byte, height int) *HW_Block {
	// 生成新的区块对象
	block := &HW_Block{
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
		height,
	}
	// 挖矿
	pow := HW_NewProofOfWork(block)
	nonce, hash := pow.HW_Run()
	block.Hw_Nonce = nonce
	block.Hw_Hash = hash[:]
	return block
}

// 将交易进行hash
func (b HW_Block) HW_HashTransactions() []byte {
	var transactions [][]byte
	// 获取交易真实内容
	for _, tx := range b.Hw_Transactions {
		transactions = append(transactions, tx.Hw_Serialize())
	}
	//txHash := sha256.Sum256(bytes.Join(transactions,[]byte{}))
	mTree := HW_NewMerkelTree(transactions)
	return mTree.Hw_RootNode.Hw_Data
}
// 新建创世区块
func HW_NewGenesisBlock(coinbase *HW_Transaction) *HW_Block {
	return HW_NewBlock([]*HW_Transaction{coinbase}, []byte{}, 1)
}
// 序列化
func (b *HW_Block) HW_Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}
// 反序列化
func HW_DeserializeBlock(d []byte) *HW_Block {
	var block HW_Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// 打印区块内容
func (block HW_Block) HWPrint() {
	fmt.Println("\n==============")
	fmt.Printf("Height:\t%d\n", block.Hw_Height)
	fmt.Printf("PrevBlockHash:\t%x\n", block.Hw_PrevBlockHash)
	fmt.Printf("Timestamp:\t%s\n", time.Unix(block.Hw_TimeStamp, 0).Format("2018-07-10 22:56:05 PM"))
	fmt.Printf("Hash:\t%x\n", block.Hw_Hash)
	fmt.Printf("Nonce:\t%d\n", block.Hw_Nonce)
	fmt.Println("Txs:")
	for _, tx := range block.Hw_Transactions {
		tx.String()
	}
	fmt.Println("区块打印结束。。。。。。。。。。。")
}
