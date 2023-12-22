package logs

import (
	"sort"
	"evm-storage-migration/utils"

	"github.com/ethereum/go-ethereum/common"
)

func Addresses(contract, topic string, addresses [][]common.Hash, ids []int, hasData bool) [][]common.Hash {
	logs := utils.Logs(common.HexToAddress(contract), common.HexToHash(topic))
	for _, vLog := range logs {
		for i, id := range ids {
			addresses[i] = append(addresses[i], vLog.Topics[id])
		}
		if hasData {
			addresses[len(ids)] = append(addresses[len(ids)], common.BytesToHash(vLog.Data))
		}
	}
	var results [][]common.Hash
	for _, array := range addresses {
		sort.Slice(array, func(i, j int) bool {
			return array[i].Big().Cmp(array[j].Big()) < 0
		})
		results = append(results, utils.RemoveDuplicate(array))
	}
	return results
}
