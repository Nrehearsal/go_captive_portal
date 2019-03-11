package wifidogkernel

import (
	"github.com/Nrehearsal/go_captive_portal/utils"
	"io/ioutil"
)

const WIFIDOG_NG = "wifidog-ng"

const INTERFACE_ARGS = "interface="
const HTTPPORT_ARGS = "port="
const HTTPSPORT_ARGS = "ssl_port="
const ENABLE_ARGS = "enabled="

const WIFIDOG_NG_PROC_FILE = "/proc/wifidog-ng/config"

func WriteToProcFile(command []byte) error {
	command = append(command, '\n')
	err := ioutil.WriteFile(WIFIDOG_NG_PROC_FILE, command, 0644)
	if err != nil {
		return err
	}
	return nil
}
func LoadModule() error {
	cmd := []string{"modprobe", WIFIDOG_NG}
	err := utils.RunCommand(cmd...)
	if err != nil {
		return err
	}
	return nil
}

func RemoveModule() error {
	cmd := []string{"rmmod", WIFIDOG_NG}
	err := utils.RunCommand(cmd...)
	if err != nil {
		return err
	}

	return nil
}

func SetGatewayInterface(interfaceName string) error {
	gatewayInterface := INTERFACE_ARGS + interfaceName
	err := WriteToProcFile([]byte(gatewayInterface))
	if err != nil {
		return err
	}
	return nil
}

func SetRedirectHttpPort(port string) error {
	httpPort := HTTPPORT_ARGS + port
	err := WriteToProcFile([]byte(httpPort))
	if err != nil {
		return err
	}
	return nil
}

func SetRedirectHttpsPort(port string) error {
	HttpsPort := HTTPSPORT_ARGS + port
	err := WriteToProcFile([]byte(HttpsPort))
	if err != nil {
		return err
	}
	return nil
}

func EnableModule() error {
	enabled := ENABLE_ARGS + "1"
	err := WriteToProcFile([]byte(enabled))
	if err != nil {
		return err
	}

	return nil
}

func DisableModule() error {
	disabled := ENABLE_ARGS + "0"
	err := WriteToProcFile([]byte(disabled))
	if err != nil {
		return err
	}

	return nil
}