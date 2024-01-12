package types

import (
	"encoding/json"
	"errors"
	"evm-storage-migration/utils"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

type Verifier struct {
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Verify  []VerifyData `json:"verify,omitempty"`
}

func (v *Verifier) TotalCases() int {
	total := 0
	for _, verify := range v.Verify {
		if len(verify.Input) > 0 {
			total += len(verify.Input)
		} else {
			total++
		}
	}
	return total
}

type VerifyData struct {
	Method string        `json:"method"`
	Input  []interface{} `json:"input,omitempty"`
	Target []string      `json:"target,omitempty"`
}

func WriteVerifier(verifier Verifier) {
	f, err := os.Create(fmt.Sprintf("./verify/%s.json", verifier.Name))
	if err != nil {
		panic(err)
	}

	defer f.Close()

	data, _ := json.MarshalIndent(verifier, "", "  ")
	_, err = f.WriteString(string(data))
	if err != nil {
		panic(err)
	}
}

func ReadVerifier(name, address string) Verifier {
	filename := fmt.Sprintf("./verify/%s.json", name)
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		WriteVerifier(Verifier{
			Name:    name,
			Address: address,
			Verify:  []VerifyData{},
		})
	}
	bytesJSON, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var verifier Verifier
	if err := json.Unmarshal(bytesJSON, &verifier); err != nil {
		panic(err)
	}
	return verifier
}

func StoreData(verifier *Verifier, verify Verify, data interface{}) {
	if i := findVerifyDataIndex(verifier, verify.Method); i != nil {
		return
	}
	if data != nil {
		input := make([]interface{}, 0)
		if verify.Addresses != nil {
			addresses := data.([]common.Hash)
			for _, address := range addresses {
				input = append(input, utils.HashToAddress(address).Hex())
			}
		} else if verify.Input.Data != nil {
			input = *verify.Input.Data
		} else {
			hash := data.(common.Hash)
			for i := 0; i < int(utils.HashToUint(hash).Int64()); i++ {
				input = append(input, i+1)
			}
		}
		verifier.Verify = append(verifier.Verify, VerifyData{
			Method: verify.Method,
			Input:  input,
			Target: verify.Target,
		})
		return
	}
	verifier.Verify = append(verifier.Verify, VerifyData{
		Method: verify.Method,
		Target: verify.Target,
	})
}

func findVerifyDataIndex(verifier *Verifier, method string) *int {
	for i, data := range verifier.Verify {
		if data.Method == method {
			return &i
		}
	}
	return nil
}
