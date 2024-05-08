package utils

import (
	"errors"
	"fmt"
	"net"
)

func GetLocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, face := range ifaces {
		if face.Flags&net.FlagUp == 0 {
			continue
		}
		if face.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrS, err := face.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrS {
			ip := getIpAddress(addr)
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("get ip address error")
}
func getIpAddress(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}

func getServerAddr(localIp string, port int) string {
	return fmt.Sprintf("%s:%d", localIp, port)
}
