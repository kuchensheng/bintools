package util

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIp(req *http.Request) string {
	ip := req.Header.Get("X-REAL-IP")
	if IP := net.ParseIP(ip); IP != nil {
		return ip
	}
	ips := req.Header.Get("X-FORWARDED-FOR")
	split := strings.Split(ips, ",")
	for _, s := range split {
		if IP := net.ParseIP(s); IP != nil {
			return s
		}
	}
	if ip1, _, err := net.SplitHostPort(req.RemoteAddr); err != nil {
		return ""
	} else if IP := net.ParseIP(ip1); IP != nil {
		return ip1
	}
	return ""
}
