package chord

import (
	"testing"
)

func TestBetweenPassCases(t *testing.T) {
	vals := [][]uint64{
		{9, 0, 10},
		{21, 15, 50},
		{max, max - 4, 3},
		{0, max, 2},
		{max, max - 1, 0},
		{max - 1, max - 2, 0},
	}
	for _, val := range vals {
		v := Between(val[0], val[1], val[2])
		if !v {
			t.Errorf("%d is between %d and %d", val[0], val[1], val[2])
		}
	}
}

func TestBetweenFailCases(t *testing.T) {
	vals := [][]uint64{
		{0, 0, 10},
		{51, 15, 50},
		{max - 4, max - 4, 3},
		{100, max, 2},
		{max - 5, max - 1, 0},
		{1000000, max - 1, 0},
	}
	for _, val := range vals {
		v := Between(val[0], val[1], val[2])
		if v {
			t.Errorf("%d is not between %d and %d", val[0], val[1], val[2])
		}
	}
}

func TestOneNode(t *testing.T) {
	n := NewNode("10000")
	n.Create()
	t.Log(n.Id)
	if n.Id > max {
		t.Errorf("This can't be > 1024!")
	}
}

func TestFindSuccessorOfOneNode(t *testing.T) {
	n := NewNode("10101")
	n.Create()
	node := n.FindSuccessor(uint64(100))
	if node != n {
		t.Errorf("The one node is its own successor")
	}
}

func TestFindSuccessorOfTwoNodes(t *testing.T) {
	n := NewNode("10101")
	n.Create()
	n2 := NewNode("10000")

	n2.Join(n)
	if n2.Fingers[0] != n {
		t.Errorf("Successor of n2 should definitely be n")
	}
}

func TestGetOneNode(t *testing.T) {
	n1 := NewNode("10000")
	n1.Create()
	n1.Data[123] = "hello world"
	data := n1.Get(123)
	if data != "hello world" {
		t.Errorf("Data should be the same")
	}
}

func TestSetTwoNodes(t *testing.T) {
	n1 := NewNode("10000")
	n1.Create()
	n2 := NewNode("10002")
	n2.Join(n1)
	n2.Stabilize()
	n1.Stabilize()
	n2.Stabilize()
	n1.Stabilize()
	t.Logf("n1: %v", n1.Id)
	t.Logf("n2: %v", n2.Id)
	t.Logf("n1 Succ: %v", n1.Fingers[0].Id)
	t.Logf("n2 Succ: %v", n2.Fingers[0].Id)
	t.Logf("n1 pred: %v", n1.Predecessor.Id)
	t.Logf("n2 pred: %v", n2.Predecessor.Id)

	n1.Set(123, "hello world")
	n1.Set(9431690416212140030, "Go on")

	if n1.Data[123] != "hello world" {
		t.Errorf("n1 should be responsible for this key")
	}
	if n1.Data[9431690416212140030] == "Go on" {
		t.Errorf("n1 Should not have this key")
	}
	if n2.Data[9431690416212140030] != "Go on" {
		t.Errorf("n2 should be responsible for this key")
	}


}


func TestStabilize(t *testing.T) {
	n1 := NewNode("10000")
	n1.Create()
	n2 := NewNode("10001")
	n3 := NewNode("10002")
	n4 := NewNode("10003")
	t.Logf("n1 ID: %v", n1.Id)
	t.Logf("n2 ID: %v", n2.Id)
	t.Logf("n3 ID: %v", n3.Id)
	t.Logf("n4 ID: %v", n4.Id)
	n2.Join(n1)
	n3.Join(n1)
	n4.Join(n1)
	t.Logf("Joined and Stabilizing")
	for i := 0; i < 5; i ++ {
		n4.Stabilize()
		n3.Stabilize()
		n2.Stabilize()
		n1.Stabilize()
	}
	t.Logf("n1 Predecessor: %v", n1.Predecessor.Id)
	t.Logf("n1 Successor: %v", n1.Fingers[0].Id)
	t.Logf("n2 Predecessor: %v", n2.Predecessor.Id)
	t.Logf("n2 Successor: %v", n2.Fingers[0].Id)
	t.Logf("n3 Predecessor: %v", n3.Predecessor.Id)
	t.Logf("n3 Successor: %v", n3.Fingers[0].Id)
	t.Logf("n4 Predecessor: %v", n4.Predecessor.Id)
	t.Logf("n4 Successor: %v", n4.Fingers[0].Id)
	if n1.Fingers[0] != n2 {
		t.Errorf("n1 successor should be n2!")
	}
	if n2.Fingers[0] != n3 {
		t.Errorf("n2 successor should be n3!")
	}
	if n3.Fingers[0] != n4 {
		t.Errorf("n3 successor should be n4!")
	}
	if n4.Fingers[0] != n1 {
		t.Errorf("n4 successor should be n1!")
	}
}
