package app

import (
	"bharvest.io/scheduled-voter/utils"
)

type Config struct {
	General struct {
		RPC                  string `toml:"rpc"`
		GRPC                 string `toml:"grpc"`
		GRPCSecureConnection bool   `toml:"grpc_secure_connection"`
	} `toml:"general"`
	Tg utils.TgConfig `toml:"tg"`
	Vote struct {
		VoteScriptPath string `toml:"vote_script_path"`
		Chain          string `toml:"chain"`
		Option         string `toml:"option"`
		Duration       int    `toml:"duration"`
	} `toml:"vote"`
}

type BaseApp struct {
	cfg *Config
	chSubProposal chan string
}
