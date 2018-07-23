package BLC

import "fmt"

func (cli *Hw_CLI) Hw_reindexUTXO(nodeID string)  {
	bc := Hw_NewBlockchain(nodeID);
	defer bc.Hw_db.Close()
	utxoset := Hw_UTXOSet{bc}
	utxoset.Hw_Reindex()
	fmt.Println("重建成功")
}
