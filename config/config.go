package config

import (
	"os"
)

const (
	RPC_URL_KEY        = "MIGRATE_FROM_URL"
	MIGRATE_TO_URL_KEY = "MIGRATE_TO_URL"
)

func GetNodeUrl() string {
	return os.Getenv(RPC_URL_KEY)
}

func GetVerifyTo() string {
	return os.Getenv(MIGRATE_TO_URL_KEY)
}
