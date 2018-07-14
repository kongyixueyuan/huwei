package BLC

func (cli *Hw_CLI) Hw_printutxo(nodeID string) {
	bc := Hw_NewBlockchain(nodeID)
	UTXOSet := Hw_UTXOSet{bc}
	defer bc.Hw_db.Close()
	UTXOSet.String()
}
