package migration

import (
	"bytes"
	"encoding/hex"
	"evm-storage-migration/config"
	"evm-storage-migration/migration/types"
	"evm-storage-migration/utils"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func Verify(target string) {
	start := time.Now()

	fromClient := utils.NewClient(config.GetNodeUrl())
	toClient := utils.NewClient(config.GetVerifyTo())

	contract := types.ReadContract(target)
	verifier := types.ReadVerifier(target)
	types.WriteVerifier(verifier)

	address := common.HexToAddress(verifier.Address)

	fmt.Println("### Verify Storage ###")
	var result bool
	for key, value := range contract.Storage {
		keyHash := common.HexToHash(key)
		valueHash := common.HexToHash(value)
		fromData := fromClient.Storage(address, keyHash)
		toData := toClient.Storage(address, keyHash)
		if fromData.Cmp(toData) != 0 {
			fmt.Printf("result not matched!!! key: %s, from: %s, to: %s\n", key, fromData.Hex(), toData.Hex())
		}
		if fromData.Cmp(valueHash) != 0 {
			fmt.Printf("result not matched!!! key: %s, out: %s, network: %s\n", key, fromData.Hex(), value)
		}
	}
	result = printResult(result)

	fmt.Println("### Verify Function ###")
	for _, verify := range verifier.Verify {
		for _, arg := range verify.Input {
			fromData := fromClient.Call(address, verify.Method, arg)
			toData := toClient.Call(address, verify.Method, arg)
			if !bytes.Equal(fromData, toData) {
				fmt.Printf("result not matched!!! method: %s, r1: %s, r2: %s\n", verify.Method, hex.EncodeToString(fromData), hex.EncodeToString(toData))
			}
		}
	}
	printResult(result)

	end := time.Now()
	diff := end.Sub(start)
	fmt.Println(diff)
}

func printResult(result bool) bool {
	if !result {
		fmt.Println("Verify OK.")
	}
	return false
}
