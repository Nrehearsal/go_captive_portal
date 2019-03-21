package authserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/Nrehearsal/go_captive_portal/config"
	"io/ioutil"
	"log"
	"net/http"
)

const HTTP_QUERY_GW_ID = "gw_id"
const HTTP_QUERY_GW_IP = "gw_address"
const HTTP_QUERY_GW_PORT = "gw_port"
const HTTP_QUERY_GW_SSL_ON = "gw_ssl_on"

const HTTP_QUERY_TOKEN = "token"

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
	PingPage   string
	LoginPage  string
	PortalPage string
	AuthPage   string
}

var serverUrl ServerUrl

type GetRequestFunc func(url string) ([]byte, error)

var doGetRequest GetRequestFunc

func Init(authServer config.AuthServer) error {
	SetServerUrl(authServer)

	if authServer.SSLOn {
		doGetRequest = httpsGetRequest
	} else {
		doGetRequest = httpGetRequest
	}

	err := CheckStatus()
	if err != nil {
		return err
	}

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
}

func CheckStatus() error {
	data, err := doGetRequest(serverUrl.PingPage)
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

	data, err := doGetRequest(authUrl)
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