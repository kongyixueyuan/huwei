package BLC

import "fmt"

func (cli *HW_CLI) HW_reindexUTXO() {
	bc := HW_NewBlockchain();
	defer bc.hw_db.Close()
	utxoset := HW_UTXOSet{bc}
	utxoset.HW_Reindex()
	fmt.Println("重建成功")
}
