package chord

import (
	"math/big"
	"testing"
)

func TestBetween(t *testing.T) {
	testTrueCases := [][]*big.Int{
		{big.NewInt(4), big.NewInt(2), big.NewInt(5)},
		{big.NewInt(5), big.NewInt(2), big.NewInt(5)},
		{big.NewInt(-10), big.NewInt(-13), big.NewInt(-8)},
	}
	testFalseCases := [][]*big.Int{
		{big.NewInt(2), big.NewInt(2), big.NewInt(5)},
		{big.NewInt(0), big.NewInt(2), big.NewInt(5)},
		{big.NewInt(6), big.NewInt(2), big.NewInt(5)},
		{big.NewInt(-1), big.NewInt(85), big.NewInt(232323)},
	}
	for _, testCase := range testTrueCases {
		result := between(testCase[0], testCase[1], testCase[2])
		if !result {
			t.Errorf("This test case should be true: %v in (%v, %v]", testCase[0], testCase[1], testCase[2])
		}
	}
	for _, testCase := range testFalseCases {
		result := between(testCase[0], testCase[1], testCase[2])
		if result {
			t.Errorf("This test case should be false: %v in (%v, %v]", testCase[0], testCase[1], testCase[2])
		}
	}
}
