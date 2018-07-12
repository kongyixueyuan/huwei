package BLC

import "log"

func (cli *HW_CLI) HW_createblockchain(address string)  {
	//验证地址是否有效
	if !HW_ValidateAddress(address){
		log.Panic("地址无效")
	}
	bc := HW_CreateBlockchain(address)
	defer bc.hw_db.Close()

	// 生成UTXOSet数据库
	UTXOSet := HW_UTXOSet{bc}
	UTXOSet.HW_Reindex()
}
