package BLC

type UTXO struct {
	TxHash []byte // Transaction 对应的 Hash
	Index  int   // Transaction 对应的 Index 索引
	Output *TXOutput
}
