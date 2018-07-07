package BLC

import "bytes"

type TXOutput struct {
	Value         int64
	Ripemd160Hash []byte //用户名
}

// 上锁
func (txoutput *TXOutput) Lock(address string) {
	publicKeyHash := Base58Decode([]byte(address))
	txoutput.Ripemd160Hash = publicKeyHash[1: len(publicKeyHash)-addressChecksumLen]
}

// 解锁
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	publicKeyHash := Base58Decode([]byte(address))
	hash160 := publicKeyHash[1: len(publicKeyHash)-addressChecksumLen]
	return bytes.Compare(hash160, txOutput.Ripemd160Hash) == 0
}

func NewTXOutput(value int64, address string) *TXOutput {
	txOutput := &TXOutput{value, nil}
	// 设置 Ripemd160Hash
	txOutput.Lock(address)
	return txOutput

}
