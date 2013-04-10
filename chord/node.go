package chord

import (
	"net/http"
	"net/rpc"
	"hash/fnv"
	"fmt"
	"net"
)

const (
	m    = uint64(64)
	max  = ^uint64(1) // ring is from 0 to 2^m -1
	null = ^uint64(0) // use max as a nil value as it is outside the valid range
)

type Data struct {
	K uint64
	V interface{}
}

type Node struct {
	Addr        string
	Id          uint64
	Data        map[uint64]interface{}
	Fingers     []*Node
	Predecessor *Node
}

func NewNode(addr string) *Node {
	hasher := fnv.New64a()
	hasher.Write([]byte(addr))
	// Number won't be bigger than max
	id := hasher.Sum64()
	node := &Node{
		Addr:        addr,
		Id:          id,
		Data:        make(map[uint64]interface{}),
		Fingers:     make([]*Node, 1),
		Predecessor: new(Node),
	}
	// A new node should look for itself as a successor until otherwise noted
	node.Predecessor.Id = null
	node.Fingers[0] = node
	go node.listenRPC()
	return node
}

func (n *Node) listenRPC() {
	rpc.Register(n)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", n.Addr))
	if e != nil {
		panic(e)
	}
	//fmt.Printf("Now listening on port %v \n", &n.Port)
	http.Serve(l, nil)
}

func (n *Node) DialRPC(port *string, reply *rpc.Client) error {
	reply, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%v", *port))
	return err
}

func (n *Node) Get(d *Data, reply *interface{}) error {
	var succ *Node
	n.FindSuccessor(&d.K, succ)
	value := succ.Data[d.K]
	reply = &value
	return nil
}

func (n *Node) Set(d *Data, reply *bool) error {
	var succ *Node
	n.FindSuccessor(&d.K, succ)
	succ.Data[d.K] = d.V
	*reply = true
	return nil
}

// Naive finder
func (n *Node) FindSuccessor(id *uint64, reply *Node) error {
	// Uh, this shouldn't recurse forever when there is only one node
	if n.Id == *id {
		reply = n
	} else if Between(*id, n.Id, (n.Fingers[0].Id+1)%max) {
		reply = n.Fingers[0]
	} else {
		// RPC recursion
		n.Fingers[0].FindSuccessor(id, reply)
	}
	return nil
}

// Setter to Join a ring
func (n *Node) Join(m *Node, reply *bool) error {
	n.Fingers[0] = m
	*reply = true
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

func (n *Node) notify(m *Node) {
	if n.Predecessor.Id == null || Between(m.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = m
	}
	if n.Fingers[0] == n {
		n.Fingers[0] = m
	}
}

// l < x <= r
func Between(x, l, r uint64) bool {
	if l > r {
		return x > l && x <= max || x >= 0 && x < r
	}
	return x > l && x <= r
}


//func (n *node) closestpreceedingnode(id int) {
//}
