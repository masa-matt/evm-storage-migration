package migration

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"evm-storage-migration/config"
	"evm-storage-migration/migration/logs"
	"evm-storage-migration/migration/storage"
	"evm-storage-migration/migration/types"
	"evm-storage-migration/utils"
)

func Migrate(target string) {
	start := time.Now()

	var contract types.Contract

	client := utils.NewClient(config.GetNodeUrl())
	settings := types.GetSettings(target)
	verifier := types.ReadVerifier(target)
	contractAddress := common.HexToAddress(settings.Address)
	contract.ContractName = settings.Name
	contract.Address = contractAddress.Hex()
	contract.Balance = client.Balance(contractAddress).String()
	contract.Nonce = strconv.FormatUint(client.Nonce(contractAddress), 10)
	contract.Bytecode = "0x" + hex.EncodeToString(client.Code(contractAddress))
	contract.Storage = map[string]string{}

	if settings.Proxy {
		data := storage.GetImplementation(client, settings.Address)
		storage.StoreStorage(&contract.Storage, data)
	}

	slr := types.GetSlither(target)

	addrs := make(map[string][][]common.Hash, 0)
	for _, setting := range settings.Addresses {
		addr := logs.Addresses(client, settings.Address, setting.Logs.Topic, types.EmptyHashArrays(setting.Logs.Size), setting.Logs.Indexes, setting.Logs.HasData)
		addrs[setting.Name] = addr
	}

	idsKey := make(map[string][]common.Hash, 0)
	for _, read := range types.FindIdsSlither(slr, settings) {
		slot := int(read.SlitherResult.Slot.Int64())
		addresses := addrs[read.Ids.Addresses][read.Ids.Index]
		for _, address := range addresses {
			// storeComment(&contract.Storage, read.SlitherResult.Name, &i)
			mappingType := types.MappingType(read.SlitherResult.TypeString)
			data := storage.GetKV(client, settings.Address, mappingType, slot, address)
			storage.StoreStorage(&contract.Storage, data)
			length := utils.HashToLength(data.Value)
			for i := 0; i < int(length.Uint64()); i++ {
				data2 := storage.GetKV(client, settings.Address, types.Mapping_key_array_uint256, slot, data.Key, i, 1)
				storage.StoreStorage(&contract.Storage, data)
				idsKey[read.SlitherResult.Name] = append(idsKey[read.SlitherResult.Name], data2.Value)
			}
		}
	}

	for _, read := range slr {
		slot := int(read.Slot.Int64())
		if strings.HasSuffix(read.TypeString, "]") {
			size := types.GetArraySize(read.TypeString, read.Value)
			for i := 0; i < size; i++ {
				// storeComment(&contract.Storage, read.Name, &i)
				structSize := types.FindStructSize(settings, read.TypeString)
				for j := 0; j < structSize; j++ {
					index := i*structSize + j
					data := storage.GetKV(client, settings.Address, types.Array, slot, index)
					storage.StoreStorage(&contract.Storage, data)
				}
			}
		} else if strings.HasPrefix(read.TypeString, "mapping") {
			mapping := types.FindTarget(settings.Mapping, read.Name)
			if mapping == nil {
				continue
			}
			switch mapping.Key {
			case "address":
				if mapping.Ids != nil {
					ids := idsKey[*mapping.Ids]
					for _, id := range ids {
						// storeComment(&contract.Storage, read.Name, &i)
						mappingType := types.MappingType(read.TypeString)
						data := storage.GetKV(client, settings.Address, mappingType, slot, id)
						storage.StoreStorage(&contract.Storage, data)
					}
				} else if mapping.Addresses != nil {
					addresses := addrs[*mapping.Addresses][*mapping.Index]
					for _, address := range addresses {
						if mapping.Struct != nil {
							// storeComment(&contract.Storage, read.Name, &i)
							structSize := mapping.Struct.Size
							for j := 0; j < structSize; j++ {
								data := storage.GetKV(client, settings.Address, types.Mapping_key_address_array_uint256, slot, address, j)
								storage.StoreStorage(&contract.Storage, data)
							}
						} else {
							// storeComment(&contract.Storage, read.Name, &i)
							mappingType := types.MappingType(read.TypeString)
							data := storage.GetKV(client, settings.Address, mappingType, slot, address)
							storage.StoreStorage(&contract.Storage, data)
							types.StoreData(&verifier, types.FindVerify(settings, read.Name), utils.HashToAddress(address).Hex())
						}
					}
				}
			case "uint256":
				if mapping.Ids != nil {
					ids := idsKey[*mapping.Ids]
					for _, id := range ids {
						// storeComment(&contract.Storage, read.Name, &i)
						data := storage.GetKV(client, settings.Address, types.Mapping_key_uint256hash, slot, id)
						storage.StoreStorage(&contract.Storage, data)
					}
				}
			}
		} else {
			// storeComment(&contract.Storage, read.Name, nil)
			data := storage.GetKV(client, settings.Address, types.Single, slot)
			storage.StoreStorage(&contract.Storage, data)
		}
	}

	types.WriteContract(contract)
	types.WriteVerifier(verifier)

	end := time.Now()
	diff := end.Sub(start)
	fmt.Println(diff)
}
