package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

// UTXO
type Transaction struct {
	// 交易hash
	TxHash []byte

	// 输入
	Vins []*TXInput

	// 输出
	Vouts []*TXOutput
}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.TxHash) == 0 && tx.Vins[0].Vout == -1

}
//Coinbase账号
// 创世区块创建时的Transaction
func NewCoinbaseTransaction(address string) *Transaction {

	// 代表消费
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}

	// 未花费交易输出
	txOutput := &TXOutput{1000, address}

	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}

	// 设置hash值
	txCoinbase.HashTransaction()

	return txCoinbase

}

// 转账时产生的Transaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {

	// 记录已花费的output
	//{hash1:[0],hash2:[2]}
	money, spendableUTXODic := blockchain.FindSpendableUTXOS(from, amount, txs)

	var txInputs []*TXInput
	var txOutputs []*TXOutput

	for txHash, indexArray := range spendableUTXODic {

		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{txHashBytes, index, from}
			txInputs = append(txInputs, txInput)
		}

	}

	// 转账

	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = &TXOutput{int64(money) - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}

	// 设置hash值
	tx.HashTransaction()

	return tx

}

func (tx *Transaction) HashTransaction() {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}
