package types

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Settings struct {
	Name      string
	Address   string
	Proxy     bool
	Struct    []Struct
	Ids       []Ids
	Mapping   []Mapping
	Addresses []Addresses
	Verify    []Verify
}

type Ids struct {
	Name      string
	Key       string
	Addresses string
	Index     int
}

type Struct struct {
	Name string
	Size int
}

type Mapping struct {
	Key       string
	Addresses *string
	Index     *int
	Ids       *string
	Target    []string
	Struct    *Struct
}

type Addresses struct {
	Name string
	Logs Logs
}

type Logs struct {
	Topic   string
	Size    int
	Indexes []int
	HasData bool
}

type Verify struct {
	Method    string
	Input     *VerifyInput
	Addresses *string
	Index     *int
	Target    []string
}

type VerifyInput struct {
	Method *string
	Type   *string
	Data   *[]interface{}
}

func (v *Verify) HasInput() bool {
	inputS := strings.Index(v.Method, "(")
	inputE := strings.Index(v.Method, ")")
	inputArray := strings.Split(v.Method[inputS+1:inputE], ",")
	return len(inputArray) != 0
}

func GetSettings(name string) Settings {
	bytesJSON, err := os.ReadFile(fmt.Sprintf("./settings/%s.json", name))
	if err != nil {
		panic(err)
	}
	var settings Settings
	if err := json.Unmarshal(bytesJSON, &settings); err != nil {
		panic(err)
	}
	return settings
}

type IdsSlither struct {
	Ids           Ids
	SlitherResult SlitherResult
}

func FindIdsSlither(slr []SlitherResult, settings Settings) []IdsSlither {
	results := make([]IdsSlither, 0)
	for _, read := range slr {
		for _, ids := range settings.Ids {
			if ids.Name == read.Name {
				results = append(results, IdsSlither{
					Ids:           ids,
					SlitherResult: read,
				})
			}
		}
	}
	return results
}

func FindTarget(mappings []Mapping, name string) *Mapping {
	for _, mapping := range mappings {
		if slices.Contains(mapping.Target, name) {
			return &mapping
		}
	}
	return nil
}

func FindStructSize(settings Settings, typeString string) int {
	start := strings.Index(typeString, "[")
	name := typeString[:start]
	for _, st := range settings.Struct {
		if st.Name == name {
			return st.Size
		}
	}
	return 0
}

func FindVerify(settings Settings, name string) *Verify {
	for _, verify := range settings.Verify {
		for _, target := range verify.Target {
			if target == name {
				return &verify
			}
		}
	}
	return nil
}
