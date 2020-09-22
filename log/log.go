package log

import (
	"log"
	"runtime"
)

func Error(err error) {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] in [%s:%d] %s; %v", fn, line, runtime.FuncForPC(pc).Name(), err)
}

func Debug(message interface{}) {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[debug] in [%s:%d] %s; %v", fn, line, runtime.FuncForPC(pc).Name(), message)
}
