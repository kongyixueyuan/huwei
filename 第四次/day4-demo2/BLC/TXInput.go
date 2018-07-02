package BLC

type TXInput struct {
	// 交易的hash
	Txhash []byte

	// 存储TXOutput在Vout里面的索引
	Vout int

	// 用户名
	ScriptSig string
}

// 判断当前的消费是否是该地址
func (txInput *TXInput) UnLockWithAddress(address string) bool {

	return txInput.ScriptSig == address

}
