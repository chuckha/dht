package main

import (
	"bufio"
	"fmt"
	"github.com/ChuckHa/go-db/chord"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

var commands = make(map[string]Command)
var runPort = 3410
var createdOrJoined = false

type Command struct {
	Name, Description string
	Run               func(args ...string)
}

func helpRun(args ...string) {
	fmt.Println("A list of available commands: ")
	fmt.Println()
	for k, v := range commands {
		fmt.Printf("%v -- %v\n", k, v.Description)
	}
	fmt.Println()
}

func portRun(args ...string) {
	if len(args) > 0 {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		runPort = port
	}
	fmt.Printf("This node will run on port %v\n", runPort)
}

func createRun(args ...string) {
	if createdOrJoined {
		fmt.Println("Cannot create after a creation or join")
		return
	}
	node := chord.NewNode()
	rpc.Register(node)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", runPort))
	if err != nil {
		panic(err)
	}
	go http.Serve(l, nil)
	createdOrJoined = true
}

func pingRun(args ...string) {
	var request string
	if len(args) > 0 {
		request = strings.Join(args, " ")
	}
	var reply string
	inter := chord.Call(fmt.Sprintf("127.0.0.1:%v", runPort), "Ping", request, reply)
	reply = inter.(string)
	fmt.Println("got back a response:", reply)
}

func putRun(args ...string) {
	var reply bool
	inter := chord.Call(fmt.Sprintf("127.0.0.1:%v", runPort), "Put", args, reply)
	reply = inter.(bool)
	if reply {
		fmt.Printf("Successfully put %v on the ring at %v\n", args[0], args[1])
	} else {
		fmt.Println("There was a problem.")
	}
}

func getRun(args ...string) {
	if len(args) != 1 {
		fmt.Println("only one argument at a time")
	}
	var reply string
	inter := chord.Call(fmt.Sprintf("127.0.0.1:%v", runPort), "Get", args[0], reply)
	reply = inter.(string)
	fmt.Println(reply)
}

func main() {
	help := Command{
		Name:        "help",
		Description: "Print a list of available commands.",
		Run:         helpRun,
	}
	quit := Command{
		Name:        "quit",
		Description: "Quit the REPL and return to the command line.",
		Run:         func(args ...string) { os.Exit(0) },
	}
	port := Command{
		Name:        "port",
		Description: fmt.Sprintf("Set the port that this node will listen on. Default: %v", runPort),
		Run:         portRun,
	}
	create := Command{
		Name:        "create",
		Description: "Create a new ring and listen on specified port",
		Run:         createRun,
	}
	ping := Command{
		Name:        "ping",
		Description: "Ping the set port.",
		Run:         pingRun,
	}
	put := Command{
		Name: "put",
		Description: "Put a key and value pair on the ring",
		Run: putRun,
	}
	get := Command{
		Name: "get",
		Description: "Get a value from a given key from the ring",
		Run: getRun,
	}

	commands[help.Name] = help
	commands[quit.Name] = quit
	commands[port.Name] = port
	commands[create.Name] = create
	commands[ping.Name] = ping
	commands[put.Name] = ping
	commands[get.Name] = ping

	for {
		fmt.Println("Chord REPL. Type ? for help")
		b := bufio.NewReader(os.Stdin)
		line, _, _ := b.ReadLine()
		args := strings.Split(string(line), " ")
		switch args[0] {
		case "help":
			help.Run()
		case "create":
			create.Run()
		case "port":
			port.Run(args[1:]...)
		case "listen":
			fmt.Println(args[1])
			fmt.Println("Listening")
		case "ping":
			ping.Run(args[1:]...)
		case "put":
			put.Run(args[1:]...)
		case "get":
			get.Run(args[1:]...)
		case "quit":
			quit.Run()
		default:
			fmt.Printf("%v: command not found. Type help to see available commands.\n", args[0])
		}
	}
}
