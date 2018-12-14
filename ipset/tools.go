package ipset

import "github.com/Nrehearsal/go_captive_portal/utils"

var WIFIDOG_NG_MAC = "wifidog-ng-mac"
var WIFIDOG_NG_IP = "wifidog-ng-ip"

var CREATE_IPSET_WIFIDOG_NG_MAC = []string{"ipset", "-!", "create", WIFIDOG_NG_MAC, "hash:mac", "timeout", "4294967"}
var CREATE_IPSET_WIFIDOG_NG_IP = []string{"ipset", "-!", "create", WIFIDOG_NG_IP, "hash:ip"}

var DESTROY_IPSET_WIFIDOG_NG_MAC = []string{"ipset", "destroy", WIFIDOG_NG_MAC}
var DESTROY_IPSET_WIFIDOG_NG_IP = []string{"ipset", "destroy", WIFIDOG_NG_IP}

var ADD_NEW_MAC_TO_IPSET = []string{"ipset", "-!", "add", WIFIDOG_NG_MAC}
var ADD_NEW_IP_TO_IPSET = []string{"ipset", "-!", "add", WIFIDOG_NG_IP}

var DEL_MAC_FROM_IPSET = []string{"ipset", "del", WIFIDOG_NG_MAC}
var DEL_IP_FROM_IPSET = []string{"ipset", "del", WIFIDOG_NG_IP}

var TEST_MAC_IN_IPSET = []string{"ipset", "test", WIFIDOG_NG_MAC}

func CreateSetForMac() error {
	err := utils.RunCommand(CREATE_IPSET_WIFIDOG_NG_MAC...)
	if err != nil {
		return err
	}
	return nil
}

func CreateSetForIp() error {
	err := utils.RunCommand(CREATE_IPSET_WIFIDOG_NG_IP...)
	if err != nil {
		return err
	}
	return nil
}

func DestroySetForMac() error {
	err := utils.RunCommand(DESTROY_IPSET_WIFIDOG_NG_MAC...)
	if err != nil {
		return err
	}
	return nil
}

func DestroySetForIp() error {
	err := utils.RunCommand(DESTROY_IPSET_WIFIDOG_NG_IP...)
	if err != nil {
		return err
	}
	return nil
}

func AddMacToSet(mac string) error {
	cmd := ADD_NEW_MAC_TO_IPSET
	cmd = append(cmd, mac)

	err := utils.RunCommand(cmd...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMacFromSet(mac string) error {
	cmd := DEL_MAC_FROM_IPSET
	cmd = append(cmd, mac)

	err := utils.RunCommand(cmd...)
	if err != nil {
		return err
	}
	return nil
}

func AddIpToSet(ip string) error {
	cmd := ADD_NEW_IP_TO_IPSET
	cmd = append(cmd, ip)

	err := utils.RunCommand(cmd...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteIpFromSet(ip string) error {
	cmd := DEL_IP_FROM_IPSET
	cmd = append(cmd, ip)

	err := utils.RunCommand(cmd...)
	if err != nil {
		return err
	}
	return nil
}

func TestMacInSet(mac string) error {
	cmd := TEST_MAC_IN_IPSET
	cmd = append(cmd, mac)

	err := utils.RunCommand(cmd...)
	if err != nil {
		return err
	}
	return nil
}