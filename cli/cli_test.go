package main

import (
	"dht"
	"testing"
)

// A call to portFn should set the Node's port
func TestPortFn(t *testing.T) {
	portFn("3333")
	if port != "3333" {
		t.Errorf("Port was not set correctly.\nExpected: %v\nGot: %v\n", "3333", port)
	}

}
