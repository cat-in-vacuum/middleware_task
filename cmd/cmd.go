package main

import (
	"context"
	"flag"
	"github.com/cat-in-vacuum/middleware_task/api"
	"github.com/cat-in-vacuum/middleware_task/limiter"
	"github.com/cat-in-vacuum/middleware_task/log"
	"github.com/cat-in-vacuum/middleware_task/notificator"
	"github.com/cat-in-vacuum/middleware_task/service"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	addr := flag.String("port", "5001", "port for app")
	flag.Parse()

	osSig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(osSig, os.Interrupt)

	ctx, cancelApi := context.WithCancel(context.Background())

	l := limiter.NewFixedWindow(limiter.Config{
		Duration: time.Minute,
		MaxReq:   100,
	})

	n := notificator.New(http.DefaultClient)
	box := service.New(n)
	server := api.New(ctx, ":"+*addr, box, l)

	go func() {
		<-osSig
		log.Debug("Server is shutting down...")
		server.Stop(context.Background(), cancelApi)
		close(done)
	}()
	server.Start()
	<-done
}
