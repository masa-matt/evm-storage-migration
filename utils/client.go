package utils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	c *ethclient.Client
}

func NewClient(url string) *Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}
	return &Client{client}
}

func (client *Client) Code(address common.Address) []byte {
	codeBytes, err := client.c.CodeAt(context.Background(), address, nil)
	if err != nil {
		panic(err)
	}
	return codeBytes
}

func (client *Client) Balance(address common.Address) *big.Int {
	balance, err := client.c.BalanceAt(context.Background(), address, nil)
	if err != nil {
		panic(err)
	}
	return balance
}

func (client *Client) Nonce(address common.Address) uint64 {
	nonce, err := client.c.NonceAt(context.Background(), address, nil)
	if err != nil {
		panic(err)
	}
	return nonce
}

func (client *Client) Logs(address common.Address, topic common.Hash) []types.Log {
	config := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}
	logs, err := client.c.FilterLogs(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return logs
}

func (client *Client) Storage(account common.Address, key common.Hash) common.Hash {
	storage, err := client.c.StorageAt(context.Background(), account, key, nil)
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(storage)
}

func (client *Client) Call(to common.Address, sig string, args ...interface{}) ([]byte, error) {
	msg := ethereum.CallMsg{
		From: to,
		To:   &to,
		Data: makeData(sig, args...),
	}
	data, err := client.c.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func makeData(sig string, args ...interface{}) []byte {
	mutability := "view"
	name, inputs, outputs, abiTypes := ParseSig(sig)
	encodedArgs := args
	if abiTypes != nil {
		encodedArgs = EncodeArgs(args, abiTypes)
	}

	method := abi.NewMethod(name, name, abi.Function, mutability, false, false, inputs, outputs)

	input, err := method.Inputs.Pack(encodedArgs...)
	if err != nil {
		panic(err)
	}

	return append(method.ID[:], input[:]...)
}
