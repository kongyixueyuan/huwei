package BLC
type Blockchain struct {
	Blocks []*Block   //存储有序的区块

}
//增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockChain(data string,height int64,preHash []byte)  {
	newBlock := NewBlock(data,height,preHash)
	blc.Blocks = append(blc.Blocks,newBlock)
}
//1. 创建带有创世区块的区块链
func CreateBlockChainWithGeneisBlock()*Blockchain{
	//创建创世区块
	genesisBlock := CreateGenesisBlock("Genenis block......")
	//返回区块链对象
	return &Blockchain{[]*Block{genesisBlock}}

}
