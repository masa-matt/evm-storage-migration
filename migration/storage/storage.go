package storage

import (
	"evm-storage-migration/migration/types"
	"evm-storage-migration/utils"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	IMPLEMENTATION_SLOT = "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
)

func StoreStorage(storage *map[string]string, data types.StorageResults) {
	(*storage)[data.Key.Hex()] = data.Value.Hex()
}

func GetImplementation(client *utils.Client, contract string) types.StorageResults {
	key := common.HexToHash(IMPLEMENTATION_SLOT)
	value := client.Storage(common.HexToAddress(contract), key)
	return types.StorageResults{
		Kvset: fmt.Sprintf("\"%s\":\"%s\"", key.Hex(), value.Hex()),
		Key:   key,
		Value: value,
	}
}

func GetKV(client *utils.Client, contract string, storageType types.StorageType, slot int, args ...interface{}) types.StorageResults {
	slotHash := common.HexToHash(fmt.Sprintf("%x", slot))
	var key common.Hash
	switch storageType {
	case types.Single:
		key = slotHash
	case types.Array:
		index := args[0].(int)
		i := new(big.Int).Add(utils.HashToUint(crypto.Keccak256Hash(slotHash.Bytes())), big.NewInt(int64(index)))
		key = common.HexToHash(fmt.Sprintf("%x", i))
	case types.Mapping_key_address:
		address := args[0].(common.Hash)
		key = crypto.Keccak256Hash(append(address.Bytes(), slotHash.Bytes()...))
	case types.Mapping_key_uint256:
		uint256 := args[0].(int)
		key = crypto.Keccak256Hash(append(common.HexToHash(fmt.Sprintf("%x", uint256)).Bytes(), slotHash.Bytes()...))
	case types.Mapping_key_uint256hash:
		uint256 := args[0].(common.Hash)
		key = crypto.Keccak256Hash(append(uint256.Bytes(), slotHash.Bytes()...))
	case types.Mapping_key_array_uint256:
		length := args[0].(common.Hash)
		index := args[1].(int)
		i := new(big.Int).Add(utils.HashToUint(crypto.Keccak256Hash(length.Bytes())), big.NewInt(int64(index)))
		key = common.HexToHash(fmt.Sprintf("%x", i))
	case types.Mapping_key_address_array_uint256:
		address := args[0].(common.Hash)
		key = crypto.Keccak256Hash(append(address.Bytes(), slotHash.Bytes()...))
		index := args[1].(int)
		i := new(big.Int).Add(utils.HashToUint(key), big.NewInt(int64(index)))
		key = common.HexToHash(fmt.Sprintf("%x", i))
	}
	value := client.Storage(common.HexToAddress(contract), key)
	return types.StorageResults{
		Kvset: fmt.Sprintf("\"%s\":\"%s\"", key.Hex(), value.Hex()),
		Key:   key,
		Value: value,
	}
}
