package BLC

import "crypto/sha256"

type Hw_MerkelTree struct {
	Hw_RootNode *Hw_MerkelNode
}

type Hw_MerkelNode struct {
	Hw_Left  *Hw_MerkelNode
	Hw_Right *Hw_MerkelNode
	Hw_Data  []byte
}

func Hw_NewMerkelTree(data [][]byte) *Hw_MerkelTree {
	var nodes []Hw_MerkelNode

	// 如果交易数据不是双数，将最后一个交易复制添加到最后
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	// 生成所有的一级节点，存储到node中
	for _, dataum := range data {
		node := Hw_NewMerkelNode(nil, nil, dataum)
		nodes = append(nodes, *node)
	}

	// 遍历生成顶层节点
	for i := 0;i<len(data)/2 ;i++{
		var newLevel []Hw_MerkelNode
		for j:=0 ; j<len(nodes) ;j+=2  {
			node := Hw_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
			newLevel = append(newLevel,*node)
		}
		nodes = newLevel
	}

	//for ; len(nodes)==1 ;{
	//	var newLevel []Hw_MerkelNode
	//	for j:=0 ; j<len(nodes) ;j+=2  {
	//		node := Hw_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
	//		newLevel = append(newLevel,*node)
	//	}
	//	nodes = newLevel
	//}
	mTree := Hw_MerkelTree{&nodes[0]}
	return &mTree
}

// 新叶节点
func Hw_NewMerkelNode(left, right *Hw_MerkelNode, data []byte) *Hw_MerkelNode {
	mNode := Hw_MerkelNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Hw_Data = hash[:]
	} else {
		prevHashes := append(left.Hw_Data, right.Hw_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Hw_Data = hash[:]
	}

	mNode.Hw_Left = left
	mNode.Hw_Right = right

	return &mNode
}
