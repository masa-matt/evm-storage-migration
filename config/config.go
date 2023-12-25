package config

import (
	"os"
)

const (
	RPC_URL_KEY = "RPC_URL"
)

func GetNodeUrl() string {
	return os.Getenv(RPC_URL_KEY)
}
