package sysinfo

import (
	"fmt"
	"net"
)

func GetMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 ||
			iface.Flags&net.FlagLoopback != 0 ||
			len(iface.HardwareAddr) == 0 {
			continue
		}

		return iface.HardwareAddr.String(), nil
	}

	return "", fmt.Errorf("no MAC address found")
}
