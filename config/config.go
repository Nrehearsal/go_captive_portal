package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type AuthServer struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	SSLOn      bool   `json:"ssl_on"`
	RootPath   string `json:"root_path"`
	PingPath   string `json:"ping_path"`
	LoginPath  string `json:"login_path"`
	PortalPath string `json:"portal_path"`
	AuthPath   string `json:"auth_path"`
}

type GatewayHttp struct {
	Port    string `json:"port"`
	SSLOn   bool   `json:"ssl_on"`
	SSLPort string `json:"ssl_port"`
	SSLCrt  string `json:"ssl_crt"`
	SSLKey  string `json:"ssl_key"`
}

type CaptiveService struct {
	GatewayInterface string      `json:"gateway_interface"`
	GWHttp           GatewayHttp `json:"gateway_http"`
	AuthServer       AuthServer  `json:"auth_server"`
}

var cpconf CaptiveService

func ParseConfigFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("读取配置文件失败，请检查配置文件是否存在: ", err.Error())
		return err
	}

	err = json.Unmarshal(data, &cpconf)
	if err != nil {
		log.Println("配置文件解析错误，请检查配置文件格式: ", err.Error())
		return err
	}
	return nil
}

func GetGatewayHttp() GatewayHttp {
	return cpconf.GWHttp
}

func GetAuthServer() AuthServer {
	return cpconf.AuthServer
}

func GetCPConf() CaptiveService {
	return cpconf
}