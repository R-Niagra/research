package queue

type Node struct {
	value int
	Next  *Node
}

func NewNode(v int) *Node {
	return &Node{
		value: v,
	}
}
