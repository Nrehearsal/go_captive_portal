package network

import (
	"errors"
	"net"
)

func DnsQueryIPv4(domain string) (string, error) {
	ips, err := net.LookupIP(domain)

	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	return "", errors.New("no Ipv4 address found")
}
