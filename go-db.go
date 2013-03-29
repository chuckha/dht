package main

import (
	"flag"
	"fmt"
	"github.com/ChuckHa/go-db/chord"
	"net"
	"net/http"
	"net/rpc"
)

var port = flag.String("port", "8000", "Port to listen on")

func main() {
	flag.Parse()

	n := chord.NewNode(*port)
	if flag.Arg(1) == "create" {
		n.Create()
	}

	registerRPC(n)
	listenRPC(port)
//	d := &chord.Data{2345, "hello world"}
//	var reply bool
//	err = client.Call("Node.Set", d, &reply)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Chord set %d: %s\n", d.K, d.V)
}

func registerRPC(n *chord.Node) {
	rpc.Register(n)
	rpc.HandleHTTP()
}

func listenRPC(port *string) {
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if e != nil {
		panic(e)
	}
	fmt.Printf("Now listening on port %v \n", *port)
	http.Serve(l, nil)
}

func dialRPC(port *string) {
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		panic(err)
	}
}

