package chord

import (
	"math/big"
)

// Node

type Node struct {
	// Strings are "address:port"
	fingers []string
	predecessor string
}

func (n *Node) findSuccessor(*big.Int) {}
func (n *Node) create() {}
func (n *Node) join(*big.Int) {}
func (n *Node) stabilize() {}
func (n *Node) notify(*big.Int) {}
func (n *Node) fixFingers() {}
func (n *Node) checkPredecessor() {}

