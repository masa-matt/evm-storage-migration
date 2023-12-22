package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func RemoveDuplicate(array []common.Hash) []common.Hash {
	keys := make(map[common.Hash]bool)
	list := []common.Hash{}

	for _, entry := range array {
		if _, value := keys[entry]; !value && entry.Cmp(common.Hash{}) != 0 {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func HashToLength(hash common.Hash) big.Int {
	return *new(big.Int).SetBytes(hash.Bytes())
}

func HashToAddress(hash common.Hash) common.Address {
	return common.BytesToAddress(hash.Bytes())
}

func HashToUint(hash common.Hash) *big.Int {
	i := new(big.Int)
	i.SetString(hash.Hex()[2:], 16)
	return i
}
