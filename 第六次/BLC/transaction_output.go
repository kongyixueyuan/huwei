package BLC

import "bytes"

type HW_TXOutput struct {
	Hw_Value      int
	Hw_PubKeyHash []byte
}

// 根据地址获取 PubKeyHash
func (out *HW_TXOutput) HW_Lock(address []byte) {
	pubKeyHash := HW_Base58Decode(address)
	pubKeyHash = pubKeyHash[1: len(pubKeyHash)-4]
	out.Hw_PubKeyHash = pubKeyHash
}

// 判断是否当前公钥对应的交易输出(是否是某个人的交易输出)
func (out *HW_TXOutput) HW_IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.Hw_PubKeyHash, pubKeyHash) == 0
}

func HW_NewTXOutput(value int, address string) *HW_TXOutput {
	txo := &HW_TXOutput{value, nil}
	txo.HW_Lock([]byte(address))
	return txo
}
