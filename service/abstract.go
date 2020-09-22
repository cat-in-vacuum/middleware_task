package service

import (
	"context"
	"github.com/cat-in-vacuum/middleware_task/notificator"
	"sync"
)


// describes an entity that provides the ability to send notifications
type Notificator interface {
	Send(ctx context.Context, notification notificator.Notification) notificator.Response
}
