package BLC

func (cli *Hw_CLI) Hw_printchain(nodeID string)  {
	Hw_NewBlockchain(nodeID).Hw_Printchain()
}
