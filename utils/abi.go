package utils

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	Uint256, _    = abi.NewType("uint256", "", nil)
	Uint32, _     = abi.NewType("uint32", "", nil)
	Uint16, _     = abi.NewType("uint16", "", nil)
	String, _     = abi.NewType("string", "", nil)
	Bool, _       = abi.NewType("bool", "", nil)
	Bytes, _      = abi.NewType("bytes", "", nil)
	Bytes32, _    = abi.NewType("bytes32", "", nil)
	Address, _    = abi.NewType("address", "", nil)
	Uint64Arr, _  = abi.NewType("uint64[]", "", nil)
	AddressArr, _ = abi.NewType("address[]", "", nil)
	Int8, _       = abi.NewType("int8", "", nil)
)

func selectType(typeString string) abi.Type {
	switch strings.TrimSpace(typeString) {
	case "uint256":
		return Uint256
	case "uint32":
		return Uint32
	case "uint16":
		return Uint16
	case "string":
		return String
	case "bool":
		return Bool
	case "bytes":
		return Bytes
	case "bytes32":
		return Bytes32
	case "address":
		return Address
	case "uint64[]":
		return Uint64Arr
	case "address[]":
		return AddressArr
	case "int8":
		return Int8
	default:
		panic(fmt.Sprintf("abi.Type not found: %s\n", typeString))
	}
}

func encodeToType(abiType abi.Type, arg interface{}) interface{} {
	switch abiType.String() {
	case Address.String():
		return common.HexToAddress(arg.(string))
	case Uint256.String(), Uint32.String(), Uint16.String(), Int8.String():
		switch a := arg.(type) {
		case string:
			b, _ := new(big.Int).SetString(a, 10)
			return b
		case float64:
			return big.NewInt(int64(a))
		default:
			return arg
		}
	default:
		return arg
	}
}

type ArgWithType struct {
	Arg     string
	AbiType abi.Type
}

func ParseSig(sig string) (string, abi.Arguments, abi.Arguments, []abi.Type) {
	inputS := strings.Index(sig, "(")
	inputE := strings.Index(sig, ")")
	outputS := strings.LastIndex(sig, "(")
	outputE := strings.Index(sig[outputS+1:], ")")
	name := sig[:inputS]
	inputArray := strings.Split(sig[inputS+1:inputE], ",")
	outputArray := strings.Split(sig[outputS+1:outputS+1+outputE], ",")
	abiTypes := make([]abi.Type, 0, len(inputArray))
	inputWithType := make([]ArgWithType, 0, len(inputArray))
	if !IsEmptyStringSlice(inputArray) {
		for _, input := range inputArray {
			abiType := selectType(input)
			abiTypes = append(abiTypes, abiType)
			inputWithType = append(inputWithType, ArgWithType{
				Arg:     input,
				AbiType: abiType,
			})
		}
	}
	outputWithType := make([]ArgWithType, 0, len(outputArray))
	if !IsEmptyStringSlice(outputArray) {
		for _, output := range outputArray {
			outputWithType = append(outputWithType, ArgWithType{
				Arg:     output,
				AbiType: selectType(output),
			})
		}
	}

	arguments := func(array []ArgWithType) []abi.Argument {
		if len(array) == 0 || array == nil || array[0].Arg == "" {
			return nil
		}
		args := make([]abi.Argument, 0, len(array))
		for _, a := range array {
			args = append(args, abi.Argument{
				Type: a.AbiType,
			})
		}
		return args
	}

	return name, arguments(inputWithType), arguments(outputWithType), abiTypes
}

func EncodeArgs(args []interface{}, abiTypes []abi.Type) []interface{} {
	encoded := make([]interface{}, 0, len(args))
	if len(abiTypes) > 1 && len(args) == 1 {
		arg, ok := args[0].(string)
		if ok {
			args = nil
			for _, a := range strings.Split(arg, ",") {
				args = append(args, strings.TrimSpace(a))
			}
		}
	}
	for i, arg := range args {
		encoded = append(encoded, encodeToType(abiTypes[i], arg))
	}
	return encoded
}
