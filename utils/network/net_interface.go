package network

import (
	"log"
	"net"
)

func GetInterfaceIP(interfaceName string) (string, error) {
	netInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", err
	}

	addrs, err := netInterface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		//Get the ipv4 address of the non-local loopback address
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}
	return "", err
}

func GetInterfaceMac(interfaceName string) (string, error) {
	netInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Println("no such a interface")
		return "", err
	}

	macAddress := netInterface.HardwareAddr
	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		log.Println("invalid hwaddr")
		return "", err
	}

	return hwAddr.String(), nil
}
