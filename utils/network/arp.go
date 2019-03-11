package network

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
)

const AWK_PATH = "awk"
const ARP_CACHE_FILE = "/proc/net/arp"

func GetMacOfIP(ip, interfaceName string) (string, error) {
	awkCondition := fmt.Sprintf(`{if ($1 == "%s" && $6 == "%s") print $4}`, ip, interfaceName)

	args := [3]string{}
	args[0] = AWK_PATH
	args[1] = awkCondition
	args[2] = ARP_CACHE_FILE

	cmd := exec.Command(args[0], args[1:]...)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}

	mac := string(data)
	if mac == "" {
		return "", errors.New("not found")
	}

	mac = mac[:len(mac)-1]
	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return "", err
	}

	return hwAddr.String(), nil
}
