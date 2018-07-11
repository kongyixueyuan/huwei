package BLC

func (cli *Rwq_CLI) HW_send(from []string, to []string, amount []string) {
	bc := HW_NewBlockchain()
	defer bc.hw_db.Close()
	bc.HWMineNewBlock(from, to, amount)
}
