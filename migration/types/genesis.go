package types

import (
	"encoding/json"
	"fmt"
	"os"
)

type Genesis struct {
	Genesis []Contract `json:"genesis"`
}

type Contract struct {
	ContractName string            `json:"contractName"`
	Balance      string            `json:"balance"`
	Nonce        string            `json:"nonce"`
	Address      string            `json:"address"`
	Bytecode     string            `json:"bytecode"`
	Storage      map[string]string `json:"storage,omitempty"`
}

func WriteContract(contract Contract) {
	f, err := os.Create(fmt.Sprintf("./out/%s.json", contract.ContractName))
	if err != nil {
		panic(err)
	}

	defer f.Close()

	data, _ := json.MarshalIndent(contract, "", "  ")
	_, err = f.WriteString(string(data))
	if err != nil {
		panic(err)
	}
}

func ReadContract(name string) Contract {
	bytesJSON, err := os.ReadFile(fmt.Sprintf("./out/%s.json", name))
	if err != nil {
		panic(err)
	}
	var contract Contract
	if err := json.Unmarshal(bytesJSON, &contract); err != nil {
		panic(err)
	}
	return contract
}
