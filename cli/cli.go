package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/chuckha/dht"
	"net/rpc"
	"os"
	"strings"
)

var (
	Node   *dht.Node
	Server *dht.Server
	port   = "3410"
)

func init() {

}

func quitFn(args ...string) error {
	Server.Quit()
	os.Exit(1)
	return nil
}
func helpFn(args ...string) error {
	fmt.Println(`Available commands:

Commands
  help         this command
  quit         exit the REPL
  port <port>  set the port for this node to listen on
  create       create a new ring with this node as the first
  join <addr>  join a ring at <addr>

Node commands
  put <key> <value>  store the value at the given key
  get <key>          get the value from a given key
  delete <key>       delete a value at a given key

Debugging commands: 
  dump         debugging output of this node
  ping <addr>  ping a node at address <addr>

`)
	return nil
}

// TODO: make sure it's a valid port
// TODO: make sure create or join hasn't been called yet
func portFn(args ...string) error {
	if Node != nil {
		return errors.New("port must be set before calling create or join")
	}
	if len(args) < 1 {
		return errors.New("Need an argument")
	}
	port = args[0]
	fmt.Printf("Port set to %v\n", args[0])
	return nil
}

func createFn(args ...string) error {
	start()
	Server.Listen()
	fmt.Printf("Node listening at %v:%v\n", Node.Address, Node.Port)
	return nil
}

func joinFn(args ...string) error {
	start()
	Server.Join(args[0])
	fmt.Printf("Joined at %v\n", args[0])
	return nil
}

func pingFn(args ...string) error {
	client, err := dial(addr())
	if err != nil {
		return err
	}
	defer client.Close()
	var response int
	err = client.Call("Node.Ping", 3, &response)
	if err != nil {
		return err
	}
	fmt.Printf("Got a response: %v\n", response)
	return nil
}

func putFn(args ...string) error {
	addr := find(args[0])
	client, err := dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	putArgs := dht.PutArgs{args[0], args[1]}
	var response bool
	err = client.Call("Node.Put", putArgs, &response)
	if err != nil {
		return err
	}
	fmt.Printf("[%v] stored %v at %v\n", addr, putArgs.Val, putArgs.Key)
	return nil
}

func getFn(args ...string) error {
	addr := find(args[0])
	client, err := dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	var response string
	err = client.Call("Node.Get", args[0], &response)
	if err != nil {
		return err
	}
	fmt.Printf("[%v] %v: %v\n", addr, args[0], response)
	return nil
}

func deleteFn(args ...string) error {
	addr := find(args[0])
	client, err := dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	var response bool
	err = client.Call("Node.Delete", args[0], &response)
	if err != nil {
		return err
	}
	fmt.Printf("[%v] deleted key %v\n", addr, args[0])
	return nil
}

func dumpFn(args ...string) error {
	fmt.Println(Server.Debug())
	return nil
}

var commands = map[string]func(args ...string) error{
	"quit":   quitFn,
	"help":   helpFn,
	"port":   portFn,
	"create": createFn,
	"delete": deleteFn,
	"ping":   pingFn,
	"put":    putFn,
	"get":    getFn,
	"dump":   dumpFn,
	"join":   joinFn,
}

func main() {
	prompt := "dht> "
	// Enter the repl
	fmt.Println(`Welcome to dht. Type "help" to learn about the available commands`)
	for {
		fmt.Printf(prompt)
		input, err := getInput()
		if err != nil {
			fmt.Println("There was an error with your command")
			continue
		}
		if _, ok := commands[input[0]]; !ok {
			fmt.Println("Command not found.")
			continue
		}
		err = commands[input[0]](input[1:]...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

/* Helpers */

func getInput() ([]string, error) {
	b := bufio.NewReader(os.Stdin)
	line, err := b.ReadString('\n')
	if err != nil {
		return []string{}, err
	}
	return strings.Split(strings.TrimSpace(line), " "), nil
}

func mustDial(addr string) *rpc.Client {
	client, err := dial(addr)
	if err != nil {
		panic(err)
	}
	return client
}

func dial(addr string) (*rpc.Client, error) {
	if addr == "" {
		addr = "127.0.0.1:" + Node.Port
	}
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func addr() string {
	return fmt.Sprintf("%v:%v", Node.Address, Node.Port)
}

// Return the address that is responsible for a given key
func find(key string) string {
	client := mustDial(addr())
	defer client.Close()
	var response string
	err := client.Call("Node.FindSuccessor", dht.Hash(key), &response)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return response
}

func start() {
	Node = dht.NewNode(port)
	Server = dht.NewServer(Node)
}
