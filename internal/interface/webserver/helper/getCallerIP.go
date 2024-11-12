package helper

import (
	"net"
	"net/http"
	"strings"
)

// getIP an attempt is made to extract the IP address from the request value
func GetIP(r *http.Request) string {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "UNKNOWN"
	}

	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	return "UNKNOWN"
}
