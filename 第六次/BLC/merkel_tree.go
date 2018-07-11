package BLC

import "crypto/sha256"

type HW_MerkelTree struct {
	Hw_RootNode *HW_MerkelNode
}

type HW_MerkelNode struct {
	Hw_Left  *HW_MerkelNode
	Hw_Right *HW_MerkelNode
	Hw_Data  []byte
}

func HW_NewMerkelTree(data [][]byte) *HW_MerkelTree {
	var nodes []HW_MerkelNode

	// 如果交易数据不是双数，将最后一个交易复制添加到最后
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	// 生成所有的一级节点，存储到node中
	for _, dataum := range data {
		node := HW_NewMerkelNode(nil, nil, dataum)
		nodes = append(nodes, *node)
	}

	// 遍历生成顶层节点
	for i := 0; i < len(data)/2; i++ {
		var newLevel []HW_MerkelNode
		for j := 0; j < len(nodes); j += 2 {
			node := HW_NewMerkelNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel
	}
	mTree := HW_MerkelTree{&nodes[0]}
	return &mTree
}

// 新叶节点
func HW_NewMerkelNode(left, right *HW_MerkelNode, data []byte) *HW_MerkelNode {
	mNode := HW_MerkelNode{}

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
