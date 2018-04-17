# DHT - Distributed Hash Table

This is an implementation of a distributed hash table as outlined in
[this paper describing a Chord][1]. Some of the material was taken from
[this assignment][2] and modified wherever I felt it needed to be modified.

For instance, some of the code provided on the project page did not feel
like idiomatic Go so I turned it into what I consider more idiotmatic code.

## Command line interface

The command line interface is the same as described in the project page.

Start the interface and type help to get a list of available commands.

    Available commands:

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

## Installation

    $ go get github.com/chuckha/dht/cli

## Usage

This project gives you a command line interface for creating a node and
joining a ring. You can start as many nodes as you like and have them join
at any other node already in the ring (or they can start thier own ring).

1. After you install the program (and assuming `$GOPATH/bin` is in your `$PATH`),
run `$ cli`.
2. You'll drop into a program that takes your input and runs your commands. Type
`help` to get a full list of commands.
3. With this first node, all we want to do is create a ring. Type `create` to start
an RPC server for this node. Type `dump` and get the value of the Address field.
4. On another computer (on the same network) or in another terminal, run the `$ cli`
command again and type `join <addr>`, but replace `<addr>` with the address
that you copied in step 3.
5. Repeat step 4 as many times as you like. You can join the ring at any node and
the nodes will correct themselves to create a ring structure.
6. Put data on the node with `put <key> <value>`. You will see the address of the
node the data got stored on. To retrieve the data, use the get command: `get <key>`.

## Misc

This project is using a linear lookup and not taking advantage of the *main* point
of the Chord paper. That means lookups are `O(n)` when they should be `O(log(n))`.

This is a fun way to play with the RPC package and also build a REPL in Go.

As always, contributions are welcome.

### License

MIT.

[1]: http://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdf "Chord: A Scalable Peer-to-peer Lookup Protocol for Internet Applications"
[2]: http://cit.cs.dixie.edu/cs/3410/asst_chord.html "Project 2: Chord Distributed Hash Table"
