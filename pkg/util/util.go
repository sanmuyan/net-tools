package util

import (
	"net"
	"strconv"
	"strings"
)

func ParsePorts(s string) (int, int) {
	var maxPort int
	var minPort int
	portsInt, err := strconv.Atoi(s)
	if err != nil {
		if strings.Contains(s, "-") {
			portRange := strings.FieldsFunc(s, func(r rune) bool {
				return r == '-'
			})
			if len(portRange) == 2 {
				minPort, err = strconv.Atoi(portRange[0])
				if err != nil {
					minPort = 0
				}
				maxPort, err = strconv.Atoi(portRange[1])
				if err != nil {
					maxPort = 0
				}
			}
		}
	} else {
		minPort = portsInt
		maxPort = portsInt
	}

	if minPort < 0 {
		minPort = 0
	}
	if maxPort < 0 {
		minPort = 0
	}
	if minPort > 65535 {
		minPort = 65535
	}
	if maxPort > 65535 {
		maxPort = 65535
	}
	if minPort > maxPort {
		maxPort = minPort
	}
	return minPort, maxPort
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func ParseIPList(s string) []string {
	// 解析传入的字符串是单个IP地址还是网段，如果是网段会解析出该网段的所有IP地址
	var ips []string
	if ip := net.ParseIP(s); ip != nil {
		ips = append(ips, ip.String())
		return ips
	}
	_, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil
	}
	minIP := ipNet.IP.Mask(ipNet.Mask)
	maxIP := ipNet.IP.Mask(ipNet.Mask)
	for i := range minIP {
		maxIP[i] |= ^ipNet.Mask[i]
	}
	for ip := minIP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
		if ip.Equal(maxIP) {
			break
		}
	}
	return ips
}
