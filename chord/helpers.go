package chord

import (
	"crypto/sha1"
	"math/big"
	"fmt"
	"net/rpc"
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

func Call(address string, method string, request interface{}, reply interface{}) (interface{}) {
	client, err := rpc.DialHTTP("tcp", address)
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = client.Call(fmt.Sprintf("Node.%v", method), &request, &reply)
	return reply
}
