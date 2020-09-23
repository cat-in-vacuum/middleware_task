package log

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

func Error(err error) {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] in [%s:%d] %s; %v", fn, line, runtime.FuncForPC(pc).Name(), err)
}

func Debug(message interface{}) {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[debug] in [%s:%d] %s; %v", fn, line, runtime.FuncForPC(pc).Name(), message)
}

func Fatal(v  ...interface{}) {
	log.Fatal(v)
}

func DebugHttpReq(r *http.Request, dur time.Duration) {
	reqPath := r.URL.Path
	method := r.Method
	addr := r.RemoteAddr
	logMsg := fmt.Sprintf("reqPath %s; reqMethod %s; addr: %s; duration: %s", reqPath, method, addr, dur.String())
	Debug(logMsg)
}
