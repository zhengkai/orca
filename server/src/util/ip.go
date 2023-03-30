package util

import (
	"encoding/binary"
	"errors"
	"net"
	"net/http"
	"strings"
)

// GetUint32IP ...
func GetUint32IP(r *http.Request) (uint32, error) {
	var ipStr string

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		ipStr = realIP
	} else {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return 0, err
		}
		ipStr = ip
	}
	ipStr = strings.TrimPrefix(ipStr, `::ffff:`)

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return 0, errors.New("Invalid IP address")
	}

	// 检查是否是IPv4
	parsedIPv4 := parsedIP.To4()
	if parsedIPv4 == nil {
		return 0, errors.New("IP address not IPv4")
	}

	// 检查是否为局域网IP
	if !parsedIP.IsPrivate() {
		return 0, errors.New("Public IP address not allowed")
	}

	// ip to uint32
	return binary.BigEndian.Uint32(parsedIPv4), nil
}

// IPString ...
func IPString(ip uint32) string {
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, ip)
	return net.IP(ipBytes).String()
}
