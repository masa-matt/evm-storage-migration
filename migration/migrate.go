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
	"evm-storage-migration/migration/report"
	"evm-storage-migration/migration/storage"
	"evm-storage-migration/migration/types"
	"evm-storage-migration/utils"
)

func Migrate(target string) {
	start := time.Now()

	var contract types.Contract

	client := utils.NewClient(config.GetNodeUrl())
	settings := types.GetSettings(target)
	verifier := types.ReadVerifier(target, settings.Address)
	reporter := report.InitGenesisReport(target)

	fmt.Println("### Create Contract Info ###")
	contractAddress := common.HexToAddress(settings.Address)
	contract.ContractName = settings.Name
	contract.Address = contractAddress.Hex()
	contract.Balance = client.Balance(contractAddress).String()
	contract.Nonce = strconv.FormatUint(client.Nonce(contractAddress), 10)
	contract.Bytecode = "0x" + hex.EncodeToString(client.Code(contractAddress))
	contract.Storage = map[string]string{}

	fmt.Println("### Get Target Addresses ###")
	slr := types.GetSlither(target)

	addressBar := utils.InitBar(len(settings.Addresses))
	addressBar.Begin()
	addrs := make(map[string][][]common.Hash, 0)
	for _, setting := range settings.Addresses {
		addr := logs.Addresses(client, settings.Address, setting.Logs.Topic, types.EmptyHashArrays(setting.Logs.Size), setting.Logs.Indexes, setting.Logs.HasData)
		addrs[setting.Name] = addr
		addressBar.Add()
	}
	addressBar.Finish()

	fmt.Println("### Get All Storage ###")
	if settings.Proxy {
		data := storage.GetImplementation(client, settings.Address)
		storage.StoreStorage(&contract.Storage, data)
		reporter.AddGenesisKV(storage.IMPLEMENTATION_NAME, storage.IMPLEMENTATION_TYPE, data.Report, data.Key.Hex(), data.Value.Hex())
	}

	idsKey := make(map[string][]common.Hash, 0)
	idsSlither := types.FindIdsSlither(slr, settings)
	storageBar := utils.InitBar(len(idsSlither) + len(slr))
	storageBar.Begin()
	for _, read := range idsSlither {
		slot := int(read.SlitherResult.Slot.Int64())
		addresses := addrs[read.Ids.Addresses][read.Ids.Index]
		for _, address := range addresses {
			mappingType := types.MappingType(read.SlitherResult.TypeString)
			data := storage.GetKV(client, settings.Address, mappingType, slot, address)
			storage.StoreStorage(&contract.Storage, data)
			reporter.AddGenesisKV(read.SlitherResult.Name, read.SlitherResult.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
			length := utils.HashToLength(data.Value)
			for i := 0; i < int(length.Uint64()); i++ {
				data2 := storage.GetKV(client, settings.Address, types.Mapping_key_array_uint256, slot, data.Key, i, 1)
				storage.StoreStorage(&contract.Storage, data)
				idsKey[read.SlitherResult.Name] = append(idsKey[read.SlitherResult.Name], data2.Value)
				reporter.AddGenesisKV(read.SlitherResult.Name, read.SlitherResult.TypeString, data2.Report, data2.Key.Hex(), data2.Value.Hex())
			}
		}
		storageBar.Add()
	}

	for _, read := range slr {
		slot := int(read.Slot.Int64())
		if strings.HasSuffix(read.TypeString, "]") {
			size := types.GetArraySize(read.TypeString, read.Value)
			for i := 0; i < size; i++ {
				structSize := types.FindStructSize(settings, read.TypeString)
				if structSize > 0 {
					for j := 0; j < structSize; j++ {
						index := i*structSize + j
						data := storage.GetKV(client, settings.Address, types.Array, slot, index)
						storage.StoreStorage(&contract.Storage, data)
						reporter.AddGenesisKV(read.Name, read.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
					}
				} else {
					data := storage.GetKV(client, settings.Address, types.Array, slot, i)
					storage.StoreStorage(&contract.Storage, data)
					reporter.AddGenesisKV(read.Name, read.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
				}
			}
		} else if strings.HasPrefix(read.TypeString, "mapping") {
			mapping := types.FindTarget(settings.Mapping, read.Name)
			if mapping == nil {
				storageBar.Add()
				continue
			}
			switch mapping.Key {
			case "address":
				if mapping.Ids != nil {
					ids := idsKey[*mapping.Ids]
					for _, id := range ids {
						mappingType := types.MappingType(read.TypeString)
						data := storage.GetKV(client, settings.Address, mappingType, slot, id)
						storage.StoreStorage(&contract.Storage, data)
						reporter.AddGenesisKV(read.Name, read.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
					}
				} else if mapping.Addresses != nil && mapping.Index != nil {
					addresses := addrs[*mapping.Addresses][*mapping.Index]
					for _, address := range addresses {
						if mapping.Struct != nil {
							structSize := mapping.Struct.Size
							for j := 0; j < structSize; j++ {
								data := storage.GetKV(client, settings.Address, types.Mapping_key_address_array_uint256, slot, address, j)
								storage.StoreStorage(&contract.Storage, data)
								reporter.AddGenesisKV(read.Name, read.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
							}
						} else {
							mappingType := types.MappingType(read.TypeString)
							data := storage.GetKV(client, settings.Address, mappingType, slot, address)
							storage.StoreStorage(&contract.Storage, data)
							reporter.AddGenesisKV(read.Name, read.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
						}
					}
				}
			case "uint256":
				if mapping.Ids != nil {
					ids := idsKey[*mapping.Ids]
					for _, id := range ids {
						data := storage.GetKV(client, settings.Address, types.Mapping_key_uint256hash, slot, id)
						storage.StoreStorage(&contract.Storage, data)
						reporter.AddGenesisKV(read.Name, read.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
					}
				}
			}
		} else {
			data := storage.GetKV(client, settings.Address, types.Single, slot)
			storage.StoreStorage(&contract.Storage, data)
			reporter.AddGenesisKV(read.Name, read.TypeString, data.Report, data.Key.Hex(), data.Value.Hex())
		}
		storageBar.Add()
	}
	storageBar.Finish()

	fmt.Println("### Create Verify Data ###")
	for _, verify := range settings.Verify {
		if verify.Addresses != nil {
			addresses := addrs[*verify.Addresses][*verify.Index]
			types.StoreData(&verifier, verify, addresses)
			continue
		}
		if verify.Input == nil {
			types.StoreData(&verifier, verify, nil)
			continue
		}
		if verify.Input.Method != nil {
			data, err := client.Call(contractAddress, *verify.Input.Method)
			if err != nil {
				panic(err)
			}
			types.StoreData(&verifier, verify, common.BytesToHash(data))
			continue
		}
		if verify.Input.Data != nil {
			types.StoreData(&verifier, verify, *verify.Input.Data)
		}
	}

	fmt.Println("### Create Outputs ###")
	types.WriteContract(contract)
	types.WriteVerifier(verifier)
	reporter.ReportVerifyResult()

	end := time.Now()
	diff := end.Sub(start)
	fmt.Println(diff)
}
