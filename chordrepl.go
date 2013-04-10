package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var commands = make(map[string]Command)
var runPort = 3410

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


func main() {
	help := Command{
		Name:        "help",
		Description: "Print a list of available commands.",
		Run:         helpRun,
	}
	quit := Command{
		Name:        "quit",
		Description: "Quit the REPL and return to the command line",
		Run:         func(args ...string) { os.Exit(0) },
	}
	port := Command{
		Name:        "port",
		Description: "Set the port that this node will listen on. Default: 3410",
		Run:         portRun,
	}
	create := Command{
		Name:        "create",
		Description: "Create a new ring",
		Run:         createRun,
	}

	commands[help.Name] = help
	commands[quit.Name] = quit
	commands[port.Name] = port

	for {
		fmt.Println("Chord REPL. Type ? for help")
		b := bufio.NewReader(os.Stdin)
		line, _, _ := b.ReadLine()
		args := strings.Split(string(line), " ")
		switch args[0] {
		case "help":
			help.Run()
		case "port":
			port.Run(args[1:]...)
		case "listen":
			fmt.Println(args[1])
			fmt.Println("Listening")
		case "put":
			fmt.Println("Putting something")
		case "quit":
			quit.Run()
		}
	}
}
