package main

import (
	"github.com/cat-in-vacuum/middleware_task/api"
	"github.com/cat-in-vacuum/middleware_task/notificator"
	"github.com/cat-in-vacuum/middleware_task/service"
	"log"
	"net/http"
)

func main() {
	//b := new(service.Box)

	// b.Notificator = notificator.New(&http.Client{})

	/*notifications := []service.Task{
		service.Task{
			URL: "example.com",
		},
		service.Task{
			URL: "example.com",
		},
		service.Task{
			URL: "example.com",
		},
		service.Task{
			URL: "example.com",
		},
		service.Task{
			URL: "example.com",
		},
		service.Task{
			URL: "example.com",
		},
	}*/

	// out := b.ProcessNotifications(context.TODO(), notifications)
	n := notificator.New(http.DefaultClient)
	box := service.New(n)
	server := api.New(":8080", box)
	log.Fatal(server.Start())
	// b.ProcessNotifications(context.TODO(), notifications)
}
