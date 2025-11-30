package utils

import (
	"errors"
	"fmt"
	"net"
)

const (
	UINT8SZ  int = 255
	UINT16SZ int = 65535
)

var ErrTooLarge = errors.New("Error: Integer too large")

// Checks if the string is a valid IP address
func CheckIP(str string) bool {
	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}
	return true
}

// Converts unsigned integer to uint8 if in range
func ToUInt8(num uint) (uint8, error) {
	if num <= uint(UINT8SZ) {
		return uint8(num), nil
	}

	return 0, ErrTooLarge
}

// Converts unsigned integer to uint16 if in range
func ToUInt16(num uint) (uint16, error) {
	if num <= uint(UINT16SZ) {
		return uint16(num), nil
	}

	return 0, ErrTooLarge
}

// Wraps error message
func WrapErr(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
