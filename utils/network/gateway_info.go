package network

import (
	"github.com/Nrehearsal/go_captive_portal/utils"
)

type GatewayInterfaceInfo struct {
	ID   string
	Name string
	IP   string
	MAC  string
}

var gatewayInfo GatewayInterfaceInfo

func GatewayInit(interfaceName string) ([]string, error) {
	mac, err := GetInterfaceMac(interfaceName)
	if err != nil {
		return nil, err
	}

	ip, err := GetInterfaceIP(interfaceName)
	if err != nil {
		return nil, err
	}

	gatewayInfo.ID = utils.GenerateMd5CipherString(mac + ip)
	gatewayInfo.Name = interfaceName
	gatewayInfo.MAC = mac
	gatewayInfo.IP = ip

	result := []string{gatewayInfo.ID, gatewayInfo.Name, gatewayInfo.MAC, gatewayInfo.IP}
	return result, nil
}

func GetInterfaceInfo() GatewayInterfaceInfo {
	return gatewayInfo
}
