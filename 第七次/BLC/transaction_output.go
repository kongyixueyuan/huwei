package BLC

import "bytes"

type Hw_TXOutput struct {
	Hw_Value  int
	Hw_PubKeyHash []byte
}
// 根据地址获取 PubKeyHash
func (out *Hw_TXOutput) Hw_Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.Hw_PubKeyHash = pubKeyHash
}

// 判断是否当前公钥对应的交易输出(是否是某个人的交易输出)
func (out *Hw_TXOutput) Hw_IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.Hw_PubKeyHash, pubKeyHash) == 0
}

func Hw_NewTXOutput(value int, address string) *Hw_TXOutput {
	txo := &Hw_TXOutput{value, nil}
	txo.Hw_Lock([]byte(address))
	return txo
}


