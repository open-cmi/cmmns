package mac

import (
	"net"
	"strings"
)

func ParseH3CMAC(macStr string) (string, error) {

	s := strings.ReplaceAll(macStr, "-", ".")
	hw, err := net.ParseMAC(s)
	if err != nil {
		return "", err
	}

	return hw.String(), nil
}

func ParseMAC(macStr string) (string, error) {

	s := strings.ReplaceAll(macStr, "-", ".")
	hw, err := net.ParseMAC(s)
	if err != nil {
		return "", err
	}

	return hw.String(), nil
}
