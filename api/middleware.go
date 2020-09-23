package api

import (
	"github.com/cat-in-vacuum/middleware_task/log"
	"net/http"
	"time"
)

func rateLimiter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tn := time.Now()
		next.ServeHTTP(w, r)
		reqDur := time.Since(tn)
		log.DebugHttpReq(r, reqDur)
	})
}
