package BLC

func (cli *Hw_CLI) Hw_send(from []string, to []string, amount []string,nodeID string, mineNow bool) {
	bc := Hw_NewBlockchain(nodeID)
	defer bc.Hw_db.Close()
	bc.MineNewBlock(from, to, amount,nodeID, mineNow)
}