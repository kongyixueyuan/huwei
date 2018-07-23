package BLC

import "log"

func (cli *Hw_CLI) Hw_createblockchain(address string,nodeID string)  {
	//验证地址是否有效
	if !Hw_ValidateAddress(address){
		log.Panic("地址无效")
	}
	bc := Hw_CreateBlockchain(address,nodeID)
	defer bc.Hw_db.Close()

	// 生成UTXOSet数据库
	UTXOSet := Hw_UTXOSet{bc}
	UTXOSet.Hw_Reindex()
}
