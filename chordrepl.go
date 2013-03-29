package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Run func (input string) (msg string, e error)
	Usage string
	Help string
}

func main() {
	for {
		fmt.Println("Chord REPL. Type ? for help")
		b := bufio.NewReader(os.Stdin)
		line, _, _ := b.ReadLine()
		args := strings.Split(string(line), " ")
		switch args[0] {
		case "listen":
			fmt.Println(args[1])
			fmt.Println("Listening")
		case "put":
			fmt.Println("Putting something")
		case "quit":
			return
		}
	}
}

