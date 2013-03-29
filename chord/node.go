package chord

import (
	"hash/fnv"
)

const (
	m = uint64(64)
	max = ^uint64(1) // ring is from 0 to 2^m -1
	null = ^uint64(0) // use max as a nil value as it is outside the valid range
)

type Data struct {
	K uint64
	V interface{}
}

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

func (n *Node) Get(d *Data, reply *interface{}) error {
	*reply = n.findSuccessor(&d.K).Data[d.K]
	return nil
}

func (n *Node) Set(d *Data, reply *bool) error {
	hasData := n.findSuccessor(&d.K)
	hasData.Data[d.K] = d.V
	return nil
}

// Naive finder
func (n *Node) findSuccessor(id *uint64) *Node {
	// Uh, this shouldn't recurse forever when there is only one node
	if n.Fingers[0] == n {
		return n
	}
	if Between(*id, n.Id, (n.Fingers[0].Id + 1) % max) {
		return n.Fingers[0]
	}
	return n.Fingers[0].findSuccessor(id)
}

// Predecessor is already defined as nil, simply find the Successor
func (n *Node) Join(m *Node, reply *bool) error {
	n.Fingers[0] = m.findSuccessor(&n.Id)
	return nil
}

// Call this every so often
func (n *Node) stabilize() {
	x := n.Fingers[0].Predecessor
	if x.Id != null && Between(x.Id, n.Id, n.Fingers[0].Id) {
		n.Fingers[0] = x
	}
	n.Fingers[0].notify(n)
}

func (n *Node) notify (m *Node) {
	if n.Predecessor.Id == null || Between(m.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = m
	}
	if n.Fingers[0] == n {
		n.Fingers[0] = m
	}
}

//func (n *node) closestpreceedingnode(id int) {
//}
