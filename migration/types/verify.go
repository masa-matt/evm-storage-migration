package types

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

type Verifier struct {
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Verify  []VerifyData `json:"verify"`
}

type VerifyData struct {
	Method string        `json:"method"`
	Input  []interface{} `json:"input"`
	Target []string      `json:"target"`
}

func WriteVerifier(verifier Verifier) {
	f, err := os.Create(fmt.Sprintf("./verify/%s.json", verifier.Name))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	data, _ := json.MarshalIndent(verifier, "", "  ")
	_, err = f.WriteString(string(data))
	if err != nil {
		log.Fatal(err)
	}
}

func ReadVerifier(name string) Verifier {
	bytesJSON, err := os.ReadFile(fmt.Sprintf("./verify/%s.json", name))
	if err != nil {
		panic(err)
	}
	var verifier Verifier
	if err := json.Unmarshal(bytesJSON, &verifier); err != nil {
		panic(err)
	}
	return verifier
}

func StoreData(verifier *Verifier, verify *Verify, data interface{}) {
	if verify == nil {
		return
	}
	if i := findVerifyDataIndex(verifier, verify.Method); i != nil {
		if len(verifier.Verify[*i].Target) == 0 {
			verifier.Verify[*i].Target = verify.Target
		}
		var cmpData interface{}
		switch v := data.(type) {
		case string:
			cmpData = strings.ToLower(v)
		default:
			cmpData = data
		}
		if !slices.Contains(verifier.Verify[*i].Input, cmpData) {
			verifier.Verify[*i].Input = append(verifier.Verify[*i].Input, data)
		}
		return
	}
	verifier.Verify = append(verifier.Verify, VerifyData{
		Method: verify.Method,
		Input:  []interface{}{data},
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
