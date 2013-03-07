package chord

import (
	"hash/fnv"
)

const (
	m = uint64(64)
	max = ^uint64(1) // ring is from 0 to 2^m -1
	null = ^uint64(0) // use max as a nil value as it is outside the valid range
)

type Node struct {
	Addr    string
	Id      uint64
	Data	map[uint64]interface{}
	Fingers []*Node
	Predecessor *Node
}

func (n *Node) Create() {
	n.Predecessor.Id = null
	n.Fingers[0] = n
}

func NewNode (addr string) *Node {
	hasher := fnv.New64a()
	hasher.Write([]byte(addr))
	// Number won't be bigger than max
	id := hasher.Sum64()
	node := Node{addr,
		id,
		make(map[uint64]interface{}),
		make([]*Node, 1),
		new(Node)}

	// A new node should look for itself as a successor until otherwise noted
	node.Predecessor.Id = null
	return &node
}

func Between(x, l, r uint64) bool {
	if l > r {
		return x > l && x <= max || x >= 0 && x < r
	}
	return  x > l && x <= r
}

func (n *Node) Get(key uint64) interface{} {
	return n.FindSuccessor(key).Data[key]
}

func (n *Node) Set(key uint64, value interface{}) {
	hasData := n.FindSuccessor(key)
	hasData.Data[key] = value
}

// Naive finder
func (n *Node) FindSuccessor(id uint64) *Node {
	// Uh, this shouldn't recurse forever when there is only one node
	if n.Fingers[0] == n {
		return n
	}
	if Between(id, n.Id, (n.Fingers[0].Id + 1) % max) {
		return n.Fingers[0]
	}
	return n.Fingers[0].FindSuccessor(id)
}

// Predecessor is already defined as nil, simply find the Successor
func (n *Node) Join(m *Node) {
	n.Fingers[0] = m.FindSuccessor(n.Id)
}

// Call this every so often
func (n *Node) Stabilize() {
	x := n.Fingers[0].Predecessor
	if x.Id != null && Between(x.Id, n.Id, n.Fingers[0].Id) {
		n.Fingers[0] = x
	}
	n.Fingers[0].Notify(n)
}

func (n *Node) Notify (m *Node) {
	if n.Predecessor.Id == null || Between(m.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = m
	}
	if n.Fingers[0] == n {
		n.Fingers[0] = m
	}
}

func (n *Node) ClosestPreceedingNode(id int) {
}
