package mac

import (
	"fmt"
	"net"
	"strings"
)

func ParseThreeSectionMAC(macStr string) (string, error) {
	if len(macStr) != 14 || strings.Count(macStr, "-") != 2 {
		return "", fmt.Errorf("invalid mac format")
	}

	s := strings.ReplaceAll(macStr, "-", ".")
	hw, err := net.ParseMAC(s)
	if err != nil {
		return "", err
	}

	return hw.String(), nil
}

func ParseMAC(macStr string) (string, error) {
	hw, err := net.ParseMAC(macStr)
	if err != nil {
		return "", err
	}

	return hw.String(), nil
}
