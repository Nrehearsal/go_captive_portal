package ipset

import (
	"go_captive_portal/config"
	"go_captive_portal/utils/network"
	"log"
	"net"
)

func Init() error {
	err := CreateSetForMac()
	if err != nil {
		log.Println("创建wifidog-ng-mac失败: ", err)
		return err
	}

	err = CreateSetForIp()
	if err != nil {
		log.Println("创建wifidog-ng-ip失败: ", err)
		return err
	}

	return nil
}

func InitWhiteList() error {
	authServer := config.GetAuthServer()
	var authServerIP string
	var err error

	ip := net.ParseIP(authServer.Host)
	if ip != nil {
		authServerIP = ip.String()
		goto ADD_AUTH_SERVER_IP_TO_SET
	}

	authServerIP, err = network.DnsQueryIPv4(authServer.Host)
	if err != nil {
		return err
	}

ADD_AUTH_SERVER_IP_TO_SET:
	//add white ip list to ipset
	AddIpToSet(authServerIP)

	gwInfo := network.GetInterfaceInfo()
	AddIpToSet(gwInfo.IP)
	whiteIpList := config.GetWhiteIpList()
	for _, ip := range whiteIpList {
		AddIpToSet(ip)
	}

	return nil
}

func Clean() error {
	err := DestroySetForMac()
	if err != nil {
		return err
	}

	err = DestroySetForIp()
	if err != nil {
		return err
	}
	return nil
}
