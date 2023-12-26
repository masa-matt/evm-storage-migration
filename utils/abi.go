package utils

import (
	"fmt"
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
	switch typeString {
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
	default:
		return arg
	}
}

func ParseSig(sig string) (string, abi.Arguments, abi.Arguments, abi.Type) {
	inputS := strings.Index(sig, "(")
	inputE := strings.Index(sig, ")")
	outputS := strings.LastIndex(sig, "(")
	outputE := strings.LastIndex(sig, ")")
	name := sig[:inputS]
	inputArray := strings.Split(sig[inputS+1:inputE], ",")
	outputArray := strings.Split(sig[outputS+1:outputE], ",")
	inputAbiType := selectType(inputArray[0])
	outputAbiType := selectType(outputArray[0])

	arguments := func(array []string, abiType abi.Type) []abi.Argument {
		args := make([]abi.Argument, 0, len(array))
		for range array {
			args = append(args, abi.Argument{
				Type: abiType,
			})
		}
		return args
	}

	return name, arguments(inputArray, inputAbiType), arguments(outputArray, outputAbiType), inputAbiType
}

func EncodeArgs(args []interface{}, abiType abi.Type) []interface{} {
	encoded := make([]interface{}, 0, len(args))
	for _, arg := range args {
		encoded = append(encoded, encodeToType(abiType, arg))
	}
	return encoded
}
