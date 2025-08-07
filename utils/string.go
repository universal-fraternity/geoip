// Package utils provides a unified processing method
package utils

import (
	"strconv"
	"strings"
)

// String2Int Convert string to int
func String2Int(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// IsIPv4 Is it IPv4
func IsIPv4(s string) bool {
	return strings.Contains(s, ".")
}

// IsIPv6 Is it IPv6
func IsIPv6(s string) bool {
	return strings.Contains(s, ":")
}
