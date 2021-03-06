package service

import (
	"context"
	"github.com/cat-in-vacuum/middleware_task/notificator"
)

const errServerShutdown = "server_shutdown"

// bus for data between place of use and Processor
type Task struct {
	URL string
}

// dependency box
type Box struct {
	notificator Notificator
}

func New(notificator Notificator) *Box {
	return &Box{notificator: notificator}
}

func (b Box) ProcessNotifications(ctx context.Context, tasks []Task) []notificator.Response {
	var (
		out = make([]notificator.Response, 0, len(tasks))

		chOut = make(chan notificator.Response, len(tasks))
		chIn  = make(chan notificator.Notification, len(tasks))
	)

	notifications := taskToNotification(tasks)
	for _, notify := range notifications {
		chIn <- notify
	}
	for i := 0; i < len(notifications); i++ {
		go b.process(ctx, chIn, chOut)
	}

	for i := 0; i < len(notifications); i++ {
		out = append(out, <-chOut)
	}
	close(chIn)
	close(chOut)

	return out
}

func (b Box) process(ctx context.Context, chIn <-chan notificator.Notification, chOut chan<- notificator.Response) {
	task := <-chIn
	chOut <- b.notificator.Send(ctx, task)
}

/*func (b Box) process(ctx context.Context, chIn <-chan notificator.Notification, chOut chan<- notificator.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	task := <-chIn
	select {
	case <-ctx.Done():
		chOut <- notificator.Response{
			Error: errServerShutdown,
			URL:   task.URL,
		}
		return
	default:
		resp := b.notificator.Send(ctx, task)
		chOut <- resp
	}

}*/

func taskToNotification(t []Task) []notificator.Notification {
	out := make([]notificator.Notification, len(t))
	for i := range t {
		out[i] = notificator.Notification{URL: t[i].URL}
	}
	return out
}
