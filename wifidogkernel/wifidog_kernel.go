package wifidogkernel

import (
	"log"
)

func Init(interfaceName, httpPort, httpsPort string) error {
	err := LoadModule()
	if err != nil {
		log.Println("wifidog-ng内核模块载入失败，请检查wifidog-ng模块是否安装: ", err)
		return err
	}

	err = SetGatewayInterface(interfaceName)
	if err != nil {
		log.Println("设置网关接口失败, 请检查网关接口名称是否正确: ", err)
		return err
	}

	err = SetRedirectHttpPort(httpPort)
	if err != nil {
		log.Println("设置本地http端口失败: ", err)
		return err
	}

	err = SetRedirectHttpsPort(httpsPort)
	if err != nil {
		log.Println("设置本地https端口失败: ", err)
		return err
	}

	err = EnableModule()
	if err != nil {
		log.Println("模块功能启用失败: ", err)
		return err
	}

	return nil
}

func Clean() error {
	err := DisableModule()
	if err != nil {
		return err
	}

	err = RemoveModule()
	if err != nil {
		return err
	}

	return nil
}