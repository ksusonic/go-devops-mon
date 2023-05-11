package agent

import (
	"fmt"
	"net"
)

func getCurrentIps() ([]net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var currentIps []net.IP
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			currentIps = append(currentIps, ip)
		}
	}
	return currentIps, nil
}

func getFirstIPOfMachine() (net.IP, error) {
	addresses, err := getCurrentIps()
	if err != nil {
		return nil, err
	} else if len(addresses) == 0 {
		return nil, fmt.Errorf("no currently addresses")
	}
	return addresses[0], nil
}
