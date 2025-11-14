package utils

import (
	"encoding/json"
	"net"
	"os"
)

type Config struct {
	BufSize uint   `json:"bufsize"`
	LogDir  string `json:"logdir"`
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

// Creates logfile if not exists
func PrepareLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// opens file for reading only
func OpenLogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY, 0o777)
}

// Writes new data to a temporary file, then replaces the real log file
func AtomicWrite(logDir, prog, logPath string, data []byte) {
	tempFile, err := os.CreateTemp(logDir, "tempfile_"+prog)
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(data)
	if err != nil {
		panic(err)
	}
	tempFile.Close()

	err = os.Rename(tempFile.Name(), logPath)
	if err != nil {
		panic(err)
	}
}
