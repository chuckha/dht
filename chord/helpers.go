package chord

import (
	"crypto/sha1"
	"math/big"
)

func hashString(val string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(val))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}

// Upper end inclusive (start, end]
func between(val, start, end *big.Int) bool {
	if val.Cmp(start) <= 0 || val.Cmp(end) > 0 {
		return false
	}
	return true
}

