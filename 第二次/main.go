package main

import (
	"publicchain/part8-proof-of-work/BLC"
	"fmt"
)

func main() {
    //1.创建一个新的块
	//block := BLC.NewBlock("Genenis block", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	//genesisBlock := BLC.CreateGenesisBlock("Genenis block......")
	//fmt.Println("genesisBlock", genesisBlock)
	//创世区块

	block := BLC.NewBlock("Test",1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
    fmt.Printf("%d\nnonce ; ",block.Nonce)
	fmt.Printf("%x\n hash : ",block.Hash)
	proofOfWork := BLC.NewProofOfWork(block)
	fmt.Printf("%v\n",proofOfWork.IsValid())
	}
