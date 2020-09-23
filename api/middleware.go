package api

import (
	"errors"
	"github.com/cat-in-vacuum/middleware_task/log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Limiter interface {
	IsAllow(ip string) bool
}

func rateLimiter(l Limiter) func(next http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				isAllow bool
			)

			ipAddr, err := getIP(r)
			if err != nil {
				isAllow = true
			} else {
				isAllow = l.IsAllow(ipAddr)
				if isAllow {
					h.ServeHTTP(w, r)
				}
			}

			w.WriteHeader(http.StatusTooManyRequests)
		})
	}
}

// https://golangbyexample.com/golang-ip-address-http-request/
func getIP(r *http.Request) (string, error) {
	// Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	// Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", errors.New("No valid ip found")
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tn := time.Now()
		next.ServeHTTP(w, r)
		reqDur := time.Since(tn)
		log.DebugHttpReq(r, reqDur)
	})
}
