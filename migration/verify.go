package migration

import (
	"evm-storage-migration/config"
	"evm-storage-migration/migration/report"
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
	vStorage := report.InitVerifyReport(target + "-storage")
	vStorageBar := utils.InitBar(len(contract.Storage))
	vStorageBar.Begin()
	for key, value := range contract.Storage {
		keyHash := common.HexToHash(key)
		fromData := fromClient.Storage(address, keyHash)
		toData := toClient.Storage(address, keyHash)
		vStorage.AddStorageResult(key, value, fromData, toData)
		vStorageBar.Add()
	}
	vStorage.ReportVerifyResult()
	vStorageBar.Finish()

	fmt.Println("### Verify Function ###")
	vFunction := report.InitVerifyReport(target + "-function")
	vFunctionBar := utils.InitBar(verifier.TotalCases())
	vFunctionBar.Finish()
	for _, verify := range verifier.Verify {
		fmt.Printf("verifying: %s\n", verify.Method)
		if len(verify.Input) > 0 {
			for _, args := range verify.Input {
				fromData := fromClient.Call(address, verify.Method, args)
				toData := toClient.Call(address, verify.Method, args)
				vFunction.AddFunctionResult(verify.Method, args, fromData, toData)
				vFunctionBar.Add()
			}
			continue
		}

		fromData := fromClient.Call(address, verify.Method)
		toData := toClient.Call(address, verify.Method)
		vFunction.AddFunctionResult(verify.Method, nil, fromData, toData)
		vFunctionBar.Add()
	}
	vFunction.ReportVerifyResult()
	vFunctionBar.Finish()

	end := time.Now()
	diff := end.Sub(start)
	fmt.Println(diff)
}
