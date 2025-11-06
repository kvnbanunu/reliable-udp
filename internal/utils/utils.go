package utils

import (
	"encoding/json"
	"net"
	"os"
)

type Config struct {
	BufSize uint   `json:"bufsize"`
	LogPath string `json:"logpath"`
	LogName string `json:"logname"`
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

func CheckIP(str string) bool {
	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}
	return true
}

// Checks if the port is within uint16 range
func CheckPort(port uint) bool {
	if port > 65535 {
		return false
	}
	return true
}
