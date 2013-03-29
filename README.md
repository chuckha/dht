# Go DB

## command line interface: `go-db`

* `go-db create -p <listening port>` creates a new ring and listens on listening port
* `go-db join -p <listening port> <addr>` connects to a node at addr and joins the ring listening on using port


## Chord

The underlying algorithm of this database is a [chord](http://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdf).
