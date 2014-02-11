package dht

import (
	"crypto/sha1"
	"math/big"
	"net"
)

// inclusive id ε (left, right]
func InclusiveBetween(left, id, right *big.Int) bool {
	// if the right is bigger than the left then we know it doesn't cross zero
	if right.Cmp(left) == 1 {
		return left.Cmp(id) == -1 && id.Cmp(right) <= 0
	}
	return left.Cmp(id) == -1 || id.Cmp(right) <= 0
}

// inclusive id ε (left, right)
func ExclusiveBetween(left, id, right *big.Int) bool {
	if right.Cmp(left) == 1 {
		return left.Cmp(id) == -1 && id.Cmp(right) < 0
	}
	return left.Cmp(id) == -1 || id.Cmp(right) < 0
}

func Hash(in string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(in))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}

const keySize = sha1.Size * 8

var hashMod = new(big.Int).Exp(big.NewInt(2), big.NewInt(keySize), nil)

func FingerEntry(start string, fingerentry int) *big.Int {
	id := Hash(start)
	two := big.NewInt(2)
	exponent := big.NewInt(int64(fingerentry) - 1)
	two.Exp(two, exponent, nil)
	id.Add(id, two)
	return id.Mod(id, hashMod)
}

func GetAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, interf := range interfaces {
		flags := interf.Flags

		// get only not loopback and up interfaces
		if flags&(net.FlagLoopback|flags&net.FlagUp) == net.FlagUp {
			addrs, err := interf.Addrs()
			if err != nil {
				panic(err)
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ip4 := ipnet.IP.To4(); len(ip4) == net.IPv4len {
						return ip4.String()
					}
				}
			}
		}

	}
	return ""
}
