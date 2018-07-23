package BLC

import (
	"log"
	"fmt"
)

func (cli *Hw_CLI) Hw_listAddrsss(nodeID string)  {
	wallets,err := Hw_NewWallets(nodeID)

	if err!=nil{
		log.Panic(err)
	}
	addresses := wallets.Hw_GetAddresses()

	for _,address := range addresses{
		fmt.Printf("%s\n",address)
	}
}
