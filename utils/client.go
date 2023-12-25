package utils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"evm-storage-migration/config"
)

func client() *ethclient.Client {
	client, err := ethclient.Dial(config.GetNodeUrl())
	if err != nil {
		panic(err)
	}
	return client
}

func Code(address common.Address) []byte {
	codeBytes, err := client().CodeAt(context.Background(), address, nil)
	if err != nil {
		panic(err)
	}
	return codeBytes
}

func Balance(address common.Address) *big.Int {
	balance, err := client().BalanceAt(context.Background(), address, nil)
	if err != nil {
		panic(err)
	}
	return balance
}

func Nonce(address common.Address) uint64 {
	nonce, err := client().NonceAt(context.Background(), address, nil)
	if err != nil {
		panic(err)
	}
	return nonce
}

func Logs(address common.Address, topic common.Hash) []types.Log {
	config := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}
	logs, err := client().FilterLogs(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return logs
}

func Storage(account common.Address, key common.Hash) common.Hash {
	storage, err := client().StorageAt(context.Background(), account, key, nil)
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(storage)
}
