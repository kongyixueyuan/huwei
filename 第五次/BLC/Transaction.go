package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"crypto/rand"
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
func NewCoinbaseTransaction(address string, amount int64) *Transaction {

	// 代表消费
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	txOutput := NewTXOutput(amount, address)
	// 未花费交易输出
	//txOutput := &TXOutput{amount, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	// 设置hash值
	txCoinbase.HashTransaction()
	return txCoinbase

}

// 转账时产生的Transaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {
	wallets, _ := NewWallets()
	wallet := wallets.WalletsMap[from];
	// 记录已花费的output
	//{hash1:[0],hash2:[2]}
	money, spendableUTXODic := blockchain.FindSpendableUTXOS(from, amount, txs)

	var txInputs []*TXInput
	var txOutputs []*TXOutput
	// 获取所有的TXInput
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{txHashBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}

	}

	// 转账

	//txOutput := &TXOutput{int64(amount), to}
	txOutput := NewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	// 找零
	//txOutput = &TXOutput{int64(money) - int64(amount), from}
	txOutput = NewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}

	// 设置hash值
	tx.HashTransaction()
	// 进行签名
	blockchain.SignTransaction(tx, wallet.Privatekey)

	return tx

}

/*
 * 序列化成字节数组
 */
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
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbaseTransaction() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Txhash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PubKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PubKey = nil

		r, s, _ := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vins[inID].Signature = signature
	}
}
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.Txhash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}

// 验证交易
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inID, vin := range tx.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.Txhash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PubKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}

func (tx *Transaction) Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
