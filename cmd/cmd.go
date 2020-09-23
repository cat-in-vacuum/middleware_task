package main

import (
	"context"
	"github.com/cat-in-vacuum/middleware_task/api"
	"github.com/cat-in-vacuum/middleware_task/log"
	"github.com/cat-in-vacuum/middleware_task/notificator"
	"github.com/cat-in-vacuum/middleware_task/service"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx, cancelApi := context.WithCancel(context.Background())
	osSig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(osSig, os.Interrupt)

	n := notificator.New(http.DefaultClient)
	box := service.New(n)
	server := api.New(ctx, ":8080", box)

	go func() {
		<-osSig
		log.Debug("Server is shutting down...")
		 server.Stop(context.Background(), cancelApi)
		close(done)
	}()
	server.Start()
	<-done
}


