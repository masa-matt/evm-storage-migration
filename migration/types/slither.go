package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"evm-storage-migration/utils"

	"github.com/ethereum/go-ethereum/common"
)

type SlitherResult struct {
	Name       string
	TypeString string
	Slot       *big.Int
	Size       *big.Int
	Offset     *big.Int
	Value      interface{}
	Elems      interface{}
}

func GetSlither(name string) []SlitherResult {
	filename := fmt.Sprintf("./slither/%s.json", name)
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	bytesJSON, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytesJSON, &result); err != nil {
		panic(err)
	}

	slitherResults := make([]SlitherResult, 0, len(result))
	for _, data := range result {
		var sl SlitherResult
		for val, d := range data.(map[string]interface{}) {
			switch val {
			case "name":
				sl.Name = d.(string)
			case "type_string":
				sl.TypeString = d.(string)
			case "slot":
				sl.Slot = big.NewInt(int64(d.(float64)))
			case "size":
				sl.Size = big.NewInt(int64(d.(float64)))
			case "offset":
				sl.Offset = big.NewInt(int64(d.(float64)))
			case "value":
				sl.Value = d
			case "elems":
				sl.Elems = d
			}
		}
		slitherResults = append(slitherResults, sl)
	}

	return slitherResults
}

func MappingType(typeString string) StorageType {
	start := strings.Index(typeString, "(")
	end := strings.Index(typeString, " =>")
	keyType := typeString[start+1 : end]
	switch keyType {
	case "address":
		return Mapping_key_address
	case "uint256":
		return Mapping_key_uint256
	}
	return Mapping_key_address
}

func GetArraySize(typeString string, value interface{}) int {
	start := strings.Index(typeString, "[")
	end := strings.Index(typeString, "]")
	size, _ := strconv.Atoi(typeString[start+1 : end])

	if size > 0 {
		return size
	}

	switch value.(type) {
	case float64:
		size = int(value.(float64))
	case string:
		size = int(utils.HashToUint(common.HexToHash(value.(string))).Int64())
	}

	return size
}
