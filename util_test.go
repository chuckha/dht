package dht

import (
	"math/big"
	"testing"
)

func TestGetAddress(t *testing.T) {
	address := GetAddress()
	if address == "" {
		t.Errorf("Expecting an address, got '%v'", address)
	}
}

func TestHash(t *testing.T) {
	hash := Hash("hello world")
	expected, _ := big.NewInt(0).SetString("243667368468580896692010249115860146898325751533", 10)
	if hash.Cmp(expected) != 0 {
		t.Errorf("Expected: %v\nGot: %v\n", expected, hash)
	}
}

type BetweenTestCase struct {
	left, id, right *big.Int
	out             bool
}

func BTC(left, id, right int64, out bool) BetweenTestCase {
	return BetweenTestCase{
		left:  big.NewInt(left),
		id:    big.NewInt(id),
		right: big.NewInt(right),
		out:   out,
	}
}

func TestInclusiveBetween(t *testing.T) {
	testCases := []BetweenTestCase{
		BTC(15, 500, 834883, true),
		BTC(3323, 12, 33299494, false),
		BTC(0, 123, 123, true),
		BTC(10, 223, 224, true),
		BTC(135, 155, 154, false),
	}

	for _, test := range testCases {
		actual := InclusiveBetween(test.left, test.id, test.right)
		if actual != test.out {
			t.Errorf("Expected %v\nGot %v\n", test.out, actual)
		}
	}

}

func TestExclusiveBetween(t *testing.T) {
	testCases := []BetweenTestCase{
		BTC(15, 500, 834883, true),
		BTC(3323, 12, 33299494, false),
		BTC(0, 123, 123, false),
		BTC(10, 223, 224, true),
		BTC(135, 155, 154, false),
	}

	for _, test := range testCases {
		actual := ExclusiveBetween(test.left, test.id, test.right)
		if actual != test.out {
			t.Errorf("Expected %v\nGot %v\n", test.out, actual)
		}
	}

}
