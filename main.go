package main

import (
	"fmt"
	chord "github.com/ChuckHa/go-db/chord"
	"net/rpc"
)

var port = 8000

func main () {
	first := chord.NewNode("8000")

	second := chord.NewNode("8001")

	var node chord.Node
	var client rpc.Client
	second.DialRPC(&first.Addr, &client)
	client.Call("Node.FindSuccessor", &second.Id, &node)
	var success bool
	second.Join(&node, &success)
	fmt.Println(first)
	fmt.Println(second)
}

//func NewChord(nodes int) {
//	previous := chord.NewNode(strconv.Itoa(port))
//	go previous.ListenRPC()
//	for i := 1; i < nodes; i++ {
//		// Get the port as a string
//		port = port + i
//		s := strconv.Itoa(port)
//		// Create the new node locally
//		n := chord.NewNode(s)
//		// Tell the new node to go and listen
//		go n.ListenRPC()
//
//		// Join the ring by
//		var reply chord.Node
//		// Talk to someone on the ring and ask:
//		// Who is my successor?
//		client := n.dialRPC(current.Addr)
//		client.Call("Node.FindSuccessor", &n.Id, &reply)
//		// The reply is the node
//		// FIXME this feels weird
//		n.Join(node)
//
//		previous = n
//	}
//}
