# go_captive_portal

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

go_captive_portal是基于[wifidog-ng](https://github.com/zhaojh329/wifidog-ng)内核模块，修改和开发的新一代无线网络强制认证方案，它具备一下几个特点：

  - 安装部署简单，一个bin，一个json配置文件即可运行
  - 摒弃复杂的iptables规则，通过netfilter模块和ipset来实现相关网络数据包操作
  - 支持http、https，https需要配合相应的操作系统（windows10/macos10.14+）和浏览器（chrome）来获得更好的体验
  - 提高了安全性，解决了通过53，67端口代理绕过认证的安全问题，通过dns，dhcp服务器列表白名单实现，修改了wifidong-ng模块的相关代码
  - 添加了数据持久化的功能，需配合配合[wifi_auth](https://github.com/Nrehearsal/wifi_auth)认证服务器使用（使用sqlite实现），或者自行实现相关业务接口
  - 添加了一些实用的API，如添加用户，查看当前在线用户列表，强制用户下线...

# Get Started

  - set up authentication server，[wifi_auth](https://github.com/Nrehearsal/wifi_auth) is good
  - `git clone https://github.com/Nrehearsal/go_captive_portal.git`
  - `cd go_captive_portal/nf_module && make`
  - `cp wifidog-ng.ko /lib/modules/$(uname -r)/kernel/ && depmod -a && modprobe wifidog-ng`
  - `cd go_captive_portal && go build`


# Configuration
  - see config.json
  
# TODO

  - Add LDAP support
