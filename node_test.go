package dht

import (
	"testing"
)

var n = NewNode("3333")

func reset() {
	n.Data = make(map[string]string)
}

// Ping should return nil
func TestPing(t *testing.T) {
	reset()
	var response int
	err := n.Ping(33, &response)
	if err != nil {
		t.Errorf("this should be nil, got %v", err)
	}
	if response != 42 {
		t.Errorf("Expected %v\nGot %v\n", 42, response)
	}
}

func TestPort(t *testing.T) {
	reset()
	args := PutArgs{"foo", "bar"}
	var response bool
	err := n.Put(args, &response)
	if err != nil {
		t.Errorf("Got error putting data %v", err)
	}
	if !response {
		t.Errorf("Failed to put data")
	}
	if n.Data["foo"] != "bar" {
		t.Errorf("Expected bar at foo, didn't get that")
	}
}

func TestGet(t *testing.T) {
	reset()
	n.Data["foo"] = "baz"
	var response string
	err := n.Get("foo", &response)
	if err != nil {
		t.Errorf("Got an error Getting a value: %v", err)
	}
	if response != "baz" {
		t.Errorf("Expected %v\n Got %v", "baz", response)
	}

}

func TestDelete(t *testing.T) {
	reset()
	n.Data["foo"] = "baz"
	var response bool
	err := n.Delete("foo", &response)
	if err != nil {
		t.Errorf("Error deleting a node %v", err)
	}
	if !response {
		t.Errorf("Expected a success, got a false")
	}
	if _, ok := n.Data["foo"]; ok {
		t.Errorf("foo was found, should not be")
	}
}
