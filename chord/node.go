package chord

import (
	"crypto/sha1"
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

// Helpers

func hashString(val string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(val))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}

// Upper end inclusive (start, end]
func between(val, start, end *big.Int) bool {
	if val.Cmp(start) <= 0 || val.Cmp(end) > 0 {
		return false
	}
	return true
}
