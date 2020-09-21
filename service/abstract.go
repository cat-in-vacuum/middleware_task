package service

import (
	"context"
	"fmt"
	"github.com/cat-in-vacuum/middleware_task/notificator"
	"sync"
)

const (
	defaultWorkersNum = 10
	defaultBuffSize   = 10
)

// describes an entity that provides the ability to send notifications
type Notificator interface {
	Send(ctx context.Context, notification notificator.Notification) *notificator.Response
}

// bus for data between place of use and Processor
type Task struct {
	URL string
}

// dependency box
type Box struct {
	Notificator Notificator
}

type WorkersPool struct {
	concurrency int
	chanIn      chan notificator.Notification
	chanOut     chan *notificator.Response
	wg          *sync.WaitGroup
}

func (w WorkersPool) Process(ctx context.Context, t []Task, nf Notificator) {
	notifications := taskToNotification(t)

	for i := range notifications {
		w.chanIn <- notifications[i]
	}

	for i := 0; i < defaultWorkersNum; i++ {
		go process(ctx, nf, w.chanIn, w.chanOut, w.wg)
	}

	w.wg.Wait()
}

func newWorkersPool(n int, buffSize int) WorkersPool {
	if n <= 0 {
		n = defaultWorkersNum
	}

	if buffSize <= 0 {
		buffSize = defaultBuffSize
	}

	return WorkersPool{
		concurrency: n,
		chanIn:      make(chan notificator.Notification, buffSize),
		chanOut:     make(chan *notificator.Response, buffSize),
		wg:          &sync.WaitGroup{},
	}
}

func (b Box) ProcessNotifications(ctx context.Context, tasks []Task) []notificator.Response {
	wpool := newWorkersPool(0, 0)
	wpool.Process(ctx, tasks, b.Notificator)
	for item := range wpool.chanOut {
		fmt.Println(item)
	}
	return nil
}

func process(ctx context.Context, n Notificator, in chan notificator.Notification, out chan *notificator.Response, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for notification := range in {
		resp := n.Send(ctx, notification)
		out <- resp
	}
}

func taskToNotification(t []Task) []notificator.Notification {
	out := make([]notificator.Notification, len(t))
	for i := range t {
		out[i] = notificator.Notification{URL: t[i].URL}
	}
	return out
}
