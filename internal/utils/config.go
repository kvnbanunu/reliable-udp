package utils

import (
	"encoding/json"
	"os"
)

// Hold default settings for all programs
type Config struct {
	ServerIP           string `json:"ServerIP"`
	ServerPort         uint16 `json:"ServerPort"`
	Timeout            uint8  `json:"Timeout"`
	MaxRetries         uint8  `json:"MaxRetries"`
	UseProxy           bool   `json:"UseProxy"`
	ProxyIP            string `json:"ProxyIP"`
	ProxyPort          uint16 `json:"ProxyPort"`
	ClientDropRate     uint8  `json:"ClientDropRate"`
	ServerDropRate     uint8  `json:"ServerDropRate"`
	ClientDelayRate    uint8  `json:"ClientDelayRate"`
	ServerDelayRate    uint8  `json:"ServerDelayRate"`
	ClientDelayTimeMin uint   `json:"ClientDelayTimeMin"`
	ClientDelayTimeMax uint   `json:"ClientDelayTimeMax"`
	ServerDelayTimeMin uint   `json:"ServerDelayTimeMin"`
	ServerDelayTimeMax uint   `json:"ServerDelayTimeMax"`
}

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
