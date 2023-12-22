package types

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
