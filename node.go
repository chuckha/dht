package dht

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	m = 161 // 161 so we can use 1-based indexing
)

type Node struct {
	Address     string
	Port        string
	Id          *big.Int
	Successor   string
	Predecessor string
	Data        map[string]string
	fingers     [m]string
	next        int
}
type PutArgs struct {
	Key, Val string
}

func NewNode(port string) *Node {
	addr := GetAddress()

	return &Node{
		Address: addr,
		Port:    port,
		Id:      Hash(fmt.Sprintf("%v:%v", addr, port)),
		Data:    make(map[string]string),
	}
}

func (n *Node) addr() string {
	return fmt.Sprintf("%v:%v", n.Address, n.Port)
}

func (n *Node) Ping(one int, two *int) error {
	*two = 42
	return nil
}

func (n *Node) Put(args PutArgs, success *bool) error {
	n.Data[args.Key] = args.Val
	*success = true
	return nil
}

func (n *Node) Get(key string, response *string) error {
	*response = n.Data[key]
	return nil
}

func (n *Node) Delete(key string, response *bool) error {
	delete(n.Data, key)
	*response = true
	return nil
}

func (n *Node) FindSuccessor(id *big.Int, successor *string) error {
	// If the id I'm looking for falls between me and my successor
	// Then the data for this id will be found on my successor
	if InclusiveBetween(n.Id, id, Hash(n.Successor)) {
		*successor = n.Successor
		return nil
	}
	var err error
	*successor, err = RPCFindSuccessor(n.Successor, id)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Notify(addr string, response *bool) error {
	if n.Predecessor == "" || ExclusiveBetween(Hash(n.Predecessor), Hash(addr), n.Id) {
		n.Predecessor = addr
	}
	return nil
}

func (n *Node) GetPredecessor(none bool, addr *string) error {
	*addr = n.Predecessor
	return nil
}

func (n *Node) join(addr string) {
	n.Predecessor = ""
	// This is saying connect to who i'm trying to join and find my successor!
	addr, err := RPCFindSuccessor(addr, Hash(n.addr()))
	if err != nil {
		fmt.Printf("Error in join %v\n", err)
		return
	}
	n.Successor = addr

}

func (n *Node) stabalize() {
	// Successor's predecessor
	predecessor, err := RPCGetPredecessor(n.Successor)
	if err == nil {
		if ExclusiveBetween(n.Id, Hash(predecessor), Hash(n.Successor)) {
			n.Successor = predecessor
		}
	}
	err = RPCNotify(n.Successor, n.addr())
	if err != nil {
		fmt.Println(err)
	}
}

func (n *Node) checkPredecessor() {
	up, err := RPCHealthCheck(n.Predecessor)
	if err != nil || !up {
		n.Predecessor = ""
	}
}

func (n *Node) fixFingers() {
	n.next += 1
	if n.next > m-1 {
		n.next = 1
	}
	var resp string
	id := FingerEntry(n.addr(), n.next)
	err := n.FindSuccessor(id, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp == "" {
		return
	}
	for InclusiveBetween(n.Id, id, Hash(resp)) {
		n.fingers[n.next] = resp
		n.next += 1
		if n.next > m-1 {
			n.next = 0
			break
		}
		id = FingerEntry(n.addr(), n.next)
	}
}

func (n *Node) create() {
	n.Predecessor = ""
	n.Successor = n.addr()
	go n.stabalizeOften()
	go n.checkPredecessorOften()
	go n.fixFingersOften()
}

func (n *Node) fixFingersOften() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			n.fixFingers()
		}
	}
}

func (n *Node) checkPredecessorOften() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			n.checkPredecessor()
		}
	}
}

func (n *Node) stabalizeOften() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			n.stabalize()
		}
	}
}

func dial(addr string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		fmt.Println("error dialing", err)
		return nil
	}
	return client
}

func RPCNotify(addr, notice string) error {
	if addr == "" {
		return errors.New("Notify: rpc address was empty")
	}
	client := dial(addr)
	if client == nil {
		return errors.New("Client was nil")
	}
	defer client.Close()
	var response bool
	return client.Call("Node.Notify", notice, &response)
}

func RPCGetPredecessor(addr string) (string, error) {
	if addr == "" {
		return "", errors.New("FindPredecessor: rpc address was empty")
	}
	client := dial(addr)
	if client == nil {
		return "", errors.New("Client was nil")
	}
	defer client.Close()
	var response string
	err := client.Call("Node.GetPredecessor", false, &response)
	if err != nil {
		return "", err
	}
	if response == "" {
		return "", errors.New("Empty predecessor")
	}
	return response, nil
}

func RPCFindSuccessor(addr string, id *big.Int) (string, error) {
	if addr == "" {
		return "", errors.New("FindSuccessor: rpc address was empty")
	}
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer client.Close()
	var response string
	err = client.Call("Node.FindSuccessor", id, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func RPCHealthCheck(addr string) (bool, error) {
	if addr == "" {
		return false, errors.New("HealthCheck: rpc address was empty")
	}
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return false, err
	}
	defer client.Close()
	var response int
	err = client.Call("Node.Ping", 101, &response)
	// handle this a bit more gracefully
	if err != nil {
		return false, err
	}
	return true, nil
}

type Server struct {
	node      *Node
	listener  net.Listener
	listening bool
}

func NewServer(n *Node) *Server {
	return &Server{
		node: n,
	}
}

func (s *Server) Listen() {
	rpc.Register(s.node)
	//	n.create()
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", s.node.Port))
	if e != nil {
		panic(e)
	}
	s.node.create()
	s.listener = l
	s.listening = true
	go http.Serve(l, nil)
}

func (s *Server) Join(addr string) {
	s.Listen()
	s.node.join(addr)
}

func (s *Server) Quit() {
	s.listener.Close()
}

func (s *Server) Listening() bool {
	return s.listening
}

func (s *Server) Debug() string {
	return fmt.Sprintf(`
ID: %v
Listening: %v
Address: %v
Data: %v
Successor: %v
Predecessor: %v
Fingers: %v
`, s.node.Id, s.Listening(), s.node.addr(), s.node.Data, s.node.Successor, s.node.Predecessor, s.node.fingers[1:])

}
