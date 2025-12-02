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
	MaxLogs            uint   `json:"MaxLogs"`
}

// Load contents of config.json or default
func LoadConfig() (*Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return DefaultConfig(), err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// returns default config in case config.json does not exist
func DefaultConfig() *Config {
	return &Config{
		ServerIP:           "127.0.0.1",
		ServerPort:         8080,
		Timeout:            2,
		MaxRetries:         2,
		UseProxy:           false,
		ProxyIP:            "127.0.0.1",
		ProxyPort:          8081,
		ClientDropRate:     10,
		ServerDropRate:     5,
		ClientDelayRate:    20,
		ServerDelayRate:    15,
		ClientDelayTimeMin: 100,
		ClientDelayTimeMax: 200,
		ServerDelayTimeMin: 150,
		ServerDelayTimeMax: 300,
		MaxLogs:            10,
	}
}
