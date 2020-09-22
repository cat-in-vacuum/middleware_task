package api

import (
	"context"
	"encoding/json"
	"github.com/cat-in-vacuum/middleware_task/log"
	"github.com/cat-in-vacuum/middleware_task/service"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	pathRoot          = "/api/v1"
	pathNotifications = pathRoot + "/notifications"
)

type API struct {
	router *mux.Router
	port   string
}

func New(port string, box *service.Box) *API {
	r := mux.NewRouter()
	// r.Use(rateLimiter())
	r.HandleFunc(pathNotifications, notificationHandler(box)).Methods("POST")

	return &API{
		router: r,
		port:   port,
	}
}

func (a *API) Start() error {
	return http.ListenAndServe(a.port, a.router)
}

func notificationHandler(box *service.Box) http.HandlerFunc {
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

		resp := box.ProcessNotifications(context.TODO(), req)

		enc := json.NewEncoder(w)
		for i := range resp {
			if err := enc.Encode(resp[i]); err != nil {
				log.Error(err)
				continue
			}
		}
	}
}
