package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type StorageType int

const (
	Single StorageType = iota
	Array
	Mapping_key_address
	Mapping_key_uint256
	Mapping_key_uint256hash
	Mapping_key_array_uint256
	Mapping_key_address_array_uint256
)

type StorageResults struct {
	Kvset  string
	Key    common.Hash
	Value  common.Hash
	Report string
}
