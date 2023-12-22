package types

import "github.com/ethereum/go-ethereum/common"

func EmptyHashArray() []common.Hash {
	return make([]common.Hash, 0)
}

func EmptyHashArrays(num int) [][]common.Hash {
	var arrays [][]common.Hash
	for i := 0; i < num; i++ {
		arrays = append(arrays, EmptyHashArray())
	}
	return arrays
}
