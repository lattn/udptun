package udptun

import (
	"os"
	"encoding/json"
)

const errorMessage = "{\"message\": \"error\"}"
const bufferSize = 65535

type Config struct {
	LocalAddr  string `json:"local_addr"`
	AllowedIps []string `json:"allowed_ips"`
	TargetAddr string `json:"target_addr"`
	ErrorMsg   string `json:"error_msg"`
	BufferSize int64 `json:"buffer_size"`
}

func parseConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	if len(config.ErrorMsg) == 0 {
		config.ErrorMsg = errorMessage
	}

	if config.BufferSize <= 0 {
		config.BufferSize = bufferSize
	}

	return &config, nil
}