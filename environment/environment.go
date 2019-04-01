package environment

import (
	"github.com/Nrehearsal/go_captive_portal/config"
	"github.com/Nrehearsal/go_captive_portal/ipset"
	"github.com/Nrehearsal/go_captive_portal/wifidogkernel"
	"log"
	"github.com/Nrehearsal/go_captive_portal/authserver"
	"github.com/Nrehearsal/go_captive_portal/utils/network"
)

func Init() error {
	cpConf := config.GetCPConf()
	var err error
	err = authserver.Init(cpConf.AuthServer)
	if err != nil {
		log.Println("认证服务器初始化失败: ", err.Error())
		return err
	}
	log.Println("认证服务器配置初始化成功: ", cpConf.AuthServer)

	gwInfo, err := network.GatewayInit(cpConf.GatewayInterface)
	if err != nil {
		log.Println("网关接口初始化失败: ", err.Error())
		return err
	}
	log.Println("网关接口初始化成功: ", gwInfo)

	err = ipset.Init()
	if err != nil {
		return err
	}
	log.Println("ipset初始化成功")

	err = ipset.InitWhiteList()
	if err != nil {
		ipset.Clean()
		log.Println("ipset白名单初始化失败: ", err.Error())
		return err
	}
	log.Println("ipset白名单初始化成功")

	err = wifidogkernel.
		Init(cpConf.GatewayInterface, cpConf.GWHttp.Port, cpConf.GWHttp.SSLPort)
	if err != nil {
		return err
	}
	log.Println("wifidog-ng内核模块初始化成功")

	log.Println("运行环境初始化完毕")
	return nil
}

func Clean() {
	err := ipset.Clean()
	if err != nil {
		return
	}
	log.Printf("ipset重置完成")

	err = wifidogkernel.Clean()
	if err != nil {
		return
	}

	log.Printf("wifidog-ng内核模块已经卸载")
	return
}
