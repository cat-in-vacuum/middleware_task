package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cat-in-vacuum/middleware_task/log"
	"github.com/cat-in-vacuum/middleware_task/service"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	pathRoot          = "/api/v1"
	pathNotifications = pathRoot + "/notifications"
)

type API struct {
	router *mux.Router
	port   string
	srv    *http.Server
}

func New(ctx context.Context, port string, box *service.Box) *API {
	r := mux.NewRouter()
	r.Use(rateLimiter)
	r.Use(logger)
	r.HandleFunc(pathNotifications, notificationHandler(ctx, box)).Methods("POST")

	return &API{
		router: r,
		port:   port,
		srv: &http.Server{
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
			Addr:         port,
			Handler:      r,
		},
	}
}

func (a *API) Start() {
	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Sprintf("Could not listen on %s: %v\n", a.port, err))
	}
}

func (a *API) Stop(ctx context.Context, cancelFunc context.CancelFunc) error {
	a.srv.RegisterOnShutdown(
		func() {
			cancelFunc()
		})
	return a.srv.Shutdown(ctx)
}

func notificationHandler(ctx context.Context, box *service.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req  = make([]service.Task, 0)
			task = service.Task{}
		)

		p := json.NewDecoder(r.Body)
		for p.More() {
			err := p.Decode(&task)
			if err != nil {
				return
			}
			req = append(req, task)
		}

		resp := box.ProcessNotifications(ctx, req)

		enc := json.NewEncoder(w)
		for i := range resp {
			if err := enc.Encode(resp[i]); err != nil {
				log.Error(err)
				continue
			}
		}
	}
}
