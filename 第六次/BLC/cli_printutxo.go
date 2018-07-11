package BLC

func (cli *Rwq_CLI) HW_printutxo() {
	bc := HW_NewBlockchain()
	UTXOSet := HW_UTXOSet{bc}
	defer bc.hw_db.Close()
	UTXOSet.String()
}
