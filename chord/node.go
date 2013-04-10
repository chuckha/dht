package chord

import (
	"math/big"
	"fmt"
)

type Node struct {
	// Strings are "address:port"
	fingers     []string
	predecessor string
	data		map[*big.Int]string
}

func NewNode() *Node {
	return &Node{
		fingers: make([]string, 161),
		predecessor: "",
		data: make(map[*big.Int]string),
	}
}

func (n *Node) findSuccessor(*big.Int) {}
func (n *Node) create()                {}
func (n *Node) join(*big.Int)          {}
func (n *Node) stabilize()             {}
func (n *Node) notify(*big.Int)        {}
func (n *Node) fixFingers()            {}
func (n *Node) checkPredecessor()      {}

// RPC methods
func (n *Node) Ping(args, reply *interface{}) error {
	in := (*args).(string)
	*reply = in
	return nil
}
func (n *Node) Get(args, reply *interface{}) error {
	fmt.Println(n.data)
	a := (*args).(string)
	fmt.Println(*hashString(a))
	fmt.Println(hashString(a))
	*reply = n.data[hashString(a)]
	fmt.Println(*reply)
	return nil
}
func (n *Node) Put(args, reply *interface{}) error {
	a := (*args).([]string)
	n.data[hashString(a[0])] = a[1]
	*reply = true
	return nil
}
//func (n *Node) Delete() error {}
