package BLC

func (cli *HW_CLI) HW_printutxo() {
	bc := HW_NewBlockchain()
	UTXOSet := HW_UTXOSet{bc}
	defer bc.hw_db.Close()
	UTXOSet.String()
}
