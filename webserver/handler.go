package webserver

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"github.com/Nrehearsal/go_captive_portal/authserver"
	"github.com/Nrehearsal/go_captive_portal/config"
	"github.com/Nrehearsal/go_captive_portal/ipset"
	"github.com/Nrehearsal/go_captive_portal/utils/network"
	"github.com/Nrehearsal/wifi_auth/template"
	"github.com/gin-gonic/gin"
)

func NotFound404(c *gin.Context) {
	method := c.Request.Method
	if method != "GET" {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	foo, _ := c.Get("GatewaySSLOn")
	//clientSSLOn.type = string, isSSL.value = "yes" or "no"
	gwSSLOn, ok := foo.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}

	clientIP := c.ClientIP()

	gwInfo := network.GetInterfaceInfo()
	clientMac, err := network.GetMacOfIP(clientIP, gwInfo.Name)
	if err != nil {
		c.String(http.StatusForbidden, "Unknown Client")
		return
	}

	//如果是已经通过认证的用户，重定向到认证服务器的portal页面
	err = ipset.TestMacInSet(clientMac)
	if err == nil {
		log.Printf("Authenticated user")
		redirectUrl := authserver.FillPortalPageParam(clientMac, clientIP, "")

		c.Redirect(http.StatusFound, redirectUrl)
		return
	}

	//未认证用户，携带原url重定向到认证服务器的login页面
	var originUrl string
	originHost := c.Request.Host
	originQuery := c.Request.URL.RawPath

	//认证服务通过gw_ssl_on参数判断跳转到http还是https
	gwHttpConf := config.GetGatewayHttp()

	var gwPort string
	if gwSSLOn == "yes" {
		originUrl = "https://" + originHost + originQuery
		gwPort = gwHttpConf.SSLPort
	} else {
		originUrl = "http://" + originHost + originQuery
		gwPort = gwHttpConf.Port
	}

	originUrl = base64.StdEncoding.EncodeToString([]byte(originUrl))

	redirectUrl := authserver.FillLoginPageParam(gwInfo.ID, gwInfo.IP, gwPort, gwSSLOn, clientMac, clientIP, originUrl)

	log.Println("redirectUrl: " + redirectUrl)
	c.Redirect(http.StatusFound, redirectUrl)
	return
}

func Auth(c *gin.Context) {
	stage := c.DefaultQuery(authserver.HTTP_QUERY_STAGE, "")
	if stage == "" {
		goto BAD_REQUEST
	}

	if stage == authserver.AUTH_STAGE_LOGIN {
		clientLogin(c)
		return
	}

	/*
		if stage == authserver.AUTH_STAGE_LOGOUT {
			clientLogout(c)
			return
		}
	*/

BAD_REQUEST:
	c.String(http.StatusBadRequest, "Bad Request")
	return
}

func clientLogin(c *gin.Context) {
	foo, _ := c.Get("GatewaySSLOn")
	//clientSSLOn.type = string, isSSL.value = "yes" or "no"
	gwSSLOn, ok := foo.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}

	token := c.DefaultQuery(authserver.HTTP_QUERY_TOKEN, "")
	originUrl := c.DefaultQuery(authserver.HTTP_QUERY_CLIENT_ORIGINAL_URL, "")
	if token == "" {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	/*
		从新获取客户但的ip和mac防止中间人攻击。
	*/
	clientIP := c.ClientIP()
	gwInfo := network.GetInterfaceInfo()
	clientMac, err := network.GetMacOfIP(clientIP, gwInfo.Name)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unknown Client")
		return
	}

	code, err := authserver.VerifyToken(token, clientMac, clientIP, authserver.AUTH_STAGE_LOGIN)
	if err != nil || code <= 0 || code >= 3 {
		//token认证失败
		//重定向到认证服务器的登陆页面
		gwInfo := network.GetInterfaceInfo()
		gwHttpInfo := config.GetGatewayHttp()

		var gwPort string
		if gwSSLOn == "yes" {
			gwPort = gwHttpInfo.SSLPort
		} else {
			gwPort = gwHttpInfo.SSLPort
		}

		redirectUrl := authserver.FillLoginPageParam(gwInfo.ID, gwInfo.IP, gwPort, gwSSLOn, clientMac, clientIP, originUrl)
		c.Redirect(http.StatusFound, redirectUrl)
		return
	}

	//认证成功
	//将通过验证的mac添加到白名单
	ipset.AddMacToSet(clientMac, code)

	//重定向到认证服务器portal页面
	redirectUrl := authserver.FillPortalPageParam(clientMac, clientIP, originUrl)
	c.Redirect(http.StatusFound, redirectUrl)

	return
}

func AddUser(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key == "" {
		c.String(http.StatusForbidden, "Unauthorized")
		return
	}

	user := template.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.String(http.StatusForbidden, "Unauthorized")
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	url := authserver.FillAddUserPageParam(key)
	log.Println(url)

	resp, err := authserver.DoPostRequest(url, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unknown Error")
		return
	}

	c.String(http.StatusOK, string(resp))
	return
}

func OnlineList(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key == "" {
		c.String(http.StatusForbidden, "Unauthorized")
		return
	}
	url := authserver.FillOnlineListPageParam(key)

	resp, err := authserver.DoGetRequest(url)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unknown Error")
		return
	}

	onlineUsers := &[]template.User{}
	err = json.Unmarshal(resp, onlineUsers)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Invalid data")
		return
	}

	c.JSON(http.StatusOK, onlineUsers)
	return
}

func KickOutUser(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	if key == "" {
		c.String(http.StatusForbidden, "Unauthorized")
		return
	}
	username := c.DefaultQuery("username", "")
	if username == "" {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}
	mac := c.DefaultQuery("mac", "")
	if mac == "" {
		c.String(http.StatusBadGateway, "Bad Request")
		return
	}

	url := authserver.FillKickOutUserPageParam(key, username, mac)
	resp, err := authserver.DoPostRequest(url, nil)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Unknown Error")
		return
	}

	//Remove user from ipset
	ipset.DeleteMacFromSet(mac)

	c.String(http.StatusOK, string(resp))
	return
}

/*
func AddUserProxy() gin.HandlerFunc {
	target := "localhost:9000"
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			r := c.Request
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.Header["my-header"] = []string{r.Header.Get("my-header")}
			// Golang camelcases headers
			delete(req.Header, "My-Header")
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
*/

/*
func clientLogout(c *gin.Context) {
	foo, _ := c.Get("GatewaySSLOn")
	//clientSSLOn.type = string, isSSL.value = "yes" or "no"
	gwSSLOn, ok := foo.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}

	clientIP := c.ClientIP()
	token := c.DefaultQuery(authserver.HTTP_QUERY_TOKEN, "")
	clientMac := c.DefaultQuery(authserver.HTTP_QUERY_CLIENT_MAC, "")

	if token == "" || clientMac == "" {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	err := authserver.VerifyToken(token, clientMac, clientIP, authserver.AUTH_STAGE_LOGOUT)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	//将mac从移除白名单
	ipset.DeleteMacFromSet(clientMac)

	gwInfo := network.GetInterfaceInfo()
	gwHttpInfo := config.GetGatewayHttp()
	var gwPort string
	if gwSSLOn == "yes" {
		gwPort = gwHttpInfo.SSLPort
	} else {
		gwPort = gwHttpInfo.SSLPort
	}

	//重定向到认证服务器的登陆页面
	redirectUrl := authserver.FillLoginPageParam(gwInfo.IP, gwInfo.IP, gwPort, gwSSLOn, clientMac, clientIP, "")
	c.Redirect(http.StatusFound, redirectUrl)

	return
}
*/
