package authserver

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"
	"go_captive_portal/config"
	"go_captive_portal/ipset"
	"go_captive_portal/template"
	"time"
)

const HTTP_QUERY_GW_ID = "gw_id"
const HTTP_QUERY_GW_IP = "gw_address"
const HTTP_QUERY_GW_PORT = "gw_port"
const HTTP_QUERY_GW_SSL_ON = "gw_ssl_on"

const HTTP_QUERY_TOKEN = "token"

const HTTP_QUERY_KEY = "key"
const HTTP_QUERY_USERNAME = "username"

const HTTP_QUERY_CLIENT_IP = "ip"
const HTTP_QUERY_CLIENT_MAC = "mac"
const HTTP_QUERY_CLIENT_ORIGINAL_URL = "url"

const HTTP_QUERY_STAGE = "stage"

const AUTH_STAGE_LOGIN = "login"
const AUTH_STAGE_LOGOUT = "logout"

const AUTH_SERVER_STATUS_NORMAL = "Pong"

const AUTH_SERVER_TOKEN_INVALID = "Auth: 0"
const AUTH_SERVER_TOKEN_VALID_STANDARD_USER = "Auth: 1"
const AUTH_SERVER_TOKEN_VALID_GUEST_USER = "Auth: 2"

type ServerUrl struct {
	PingPage        string
	LoginPage       string
	PortalPage      string
	AuthPage        string
	AddUserPage     string
	OnlineListPage  string
	KickOutUserPage string
}

var serverUrl ServerUrl

type GetRequestFunc func(url string) ([]byte, error)
type PostRequestFunc func(url string, data []byte) ([]byte, error)

var DoGetRequest GetRequestFunc
var DoPostRequest PostRequestFunc

func Init(authServer config.AuthServer) error {
	SetServerUrl(authServer)

	if authServer.SSLOn {
		DoGetRequest = httpsGetRequest
		DoPostRequest = httpsPostRequest
	} else {
		DoGetRequest = httpGetRequest
		DoPostRequest = httpPostRequest
	}

	err := CheckStatus()
	if err != nil {
		return err
	}

	//RestoreOnlineUser()
	return nil
}

func httpsGetRequest(url string) ([]byte, error) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Transport: &tr,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return body, nil
}

func httpGetRequest(url string) ([]byte, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil
}

func httpsPostRequest(url string, data []byte) ([]byte, error) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Transport: &tr,
	}

	var buffer *bytes.Buffer
	if data != nil {
		buffer = bytes.NewBuffer(data)
	} else {
		buffer = nil
	}

	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return body, nil
}

func httpPostRequest(url string, data []byte) ([]byte, error) {
	client := http.Client{}

	var buffer *bytes.Buffer
	if data != nil {
		buffer = bytes.NewBuffer(data)
	} else {
		buffer = nil
	}

	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return body, nil
}

func SetServerUrl(AuthServer config.AuthServer) {
	var base string
	if AuthServer.SSLOn {
		base = "https://" + AuthServer.Host + ":" + AuthServer.Port + AuthServer.RootPath
	} else {
		base = "http://" + AuthServer.Host + ":" + AuthServer.Port + AuthServer.RootPath
	}

	serverUrl.PingPage = base + AuthServer.PingPath
	serverUrl.LoginPage = base + AuthServer.LoginPath
	serverUrl.AuthPage = base + AuthServer.AuthPath
	serverUrl.PortalPage = base + AuthServer.PortalPath
	serverUrl.AddUserPage = base + AuthServer.AddUserPath
	serverUrl.OnlineListPage = base + AuthServer.OnlineListPath
	serverUrl.KickOutUserPage = base + AuthServer.KickOutUserPath
}

func CheckStatus() error {
	data, err := DoGetRequest(serverUrl.PingPage)
	if err != nil {
		return err
	}

	if string(data) != AUTH_SERVER_STATUS_NORMAL {
		return errors.New("auth server offline: " + string(data))
	}

	return nil
}

func VerifyToken(token, clientMac, clientIP, stage string) (int, error) {
	authUrl := FillAuthPageParam(token, clientMac, clientIP, stage)
	log.Println(authUrl)

	data, err := DoGetRequest(authUrl)
	if err != nil {
		return 0, err
	}

	if string(data) == AUTH_SERVER_TOKEN_INVALID {
		return 0, errors.New("token is invalid: " + string(data))
	}

	if string(data) == AUTH_SERVER_TOKEN_VALID_STANDARD_USER {
		return 1, nil
	}

	if string(data) == AUTH_SERVER_TOKEN_VALID_GUEST_USER {
		return 2, nil
	}

	return 0, errors.New("unknown error")
}

func Update(config config.AuthServer) error {
	SetServerUrl(config)
	return nil
}

func FillLoginPageParam(gwId, gwIP, gwPort, gwSSLOn, clientMac, clientIP, originUrl string) string {
	url :=
		fmt.Sprintf(`%s?%s=%s&%s=%s&%s=%s&%s=%s&%s=%s&%s=%s&%s=%s`,
			serverUrl.LoginPage,
			HTTP_QUERY_GW_ID, gwId, HTTP_QUERY_GW_IP, gwIP, HTTP_QUERY_GW_PORT, gwPort, HTTP_QUERY_GW_SSL_ON, gwSSLOn,
			HTTP_QUERY_CLIENT_MAC, clientMac, HTTP_QUERY_CLIENT_IP, clientIP, HTTP_QUERY_CLIENT_ORIGINAL_URL, originUrl,
		)

	return url
}

func FillAuthPageParam(token, clientMac, clientIP, stage string) string {
	url := fmt.Sprintf(`%s?%s=%s&%s=%s&%s=%s&%s=%s`,
		serverUrl.AuthPage, HTTP_QUERY_TOKEN, token, HTTP_QUERY_CLIENT_MAC, clientMac,
		HTTP_QUERY_CLIENT_IP, clientIP, HTTP_QUERY_STAGE, stage,
	)

	return url
}

func FillPortalPageParam(clientMac, clientIP, originUrl string) string {
	url := fmt.Sprintf(`%s?%s=%s&%s=%s&%s=%s`,
		serverUrl.PortalPage, HTTP_QUERY_CLIENT_MAC, clientMac, HTTP_QUERY_CLIENT_IP,
		clientIP, HTTP_QUERY_CLIENT_ORIGINAL_URL, originUrl,
	)

	return url
}

func FillAddUserPageParam(key string) string {
	url := fmt.Sprintf(`%s?%s=%s`,
		serverUrl.AddUserPage, HTTP_QUERY_KEY, key,
	)
	return url
}

func FillOnlineListPageParam(key string) string {
	url := fmt.Sprintf(`%s?%s=%s`,
		serverUrl.OnlineListPage, HTTP_QUERY_KEY, key,
	)
	return url
}

func FillKickOutUserPageParam(key, username, mac string) string {
	url := fmt.Sprintf(`%s?%s=%s&%s=%s&%s=%s`,
		serverUrl.KickOutUserPage, HTTP_QUERY_KEY, key, HTTP_QUERY_USERNAME, username, HTTP_QUERY_CLIENT_MAC, mac,
	)
	return url
}

func RestoreOnlineUser() error {
	authServerConf := config.GetAuthServer()
	url := FillOnlineListPageParam(authServerConf.Key)

	resp, err := DoGetRequest(url)
	if err != nil {
		return err
	}

	onlineUsers := &[]template.OnlineUser{}
	err = json.Unmarshal(resp, onlineUsers)
	if err != nil {
		return err
	}

	for _, v := range *onlineUsers {
		leftTime := v.ExpiredTimeStamp - time.Now().Unix()
		if leftTime <= 0 {
			continue
		}
		log.Println(v.Mac, leftTime)
		ipset.AddMacToSetWithTimeout(v.Mac, leftTime)
	}

	return nil
}
