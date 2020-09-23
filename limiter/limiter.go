package limiter

import (
	"fmt"
	"sync"
	"time"
)

type FixedWindow struct {
	cgf     Config
	mu      sync.Mutex
	clients map[string]*counter
}

type counter struct {
	current   int
	timeStart time.Time
	timeEnd   time.Time
	resetFunc func(c *counter, t *time.Timer)
}

func (c *counter) run(dur time.Duration) {
	go c.resetFunc(c, time.NewTimer(dur))
}

func (c *counter) reset() {
	c = new(counter)
}

type Config struct {
	Duration time.Duration
	MaxReq   int
}

func NewFixedWindow(cfg Config) *FixedWindow {
	return &FixedWindow{
		cgf:     cfg,
		mu:      sync.Mutex{},
		clients: make(map[string]*counter),
	}
}

//
func (l *FixedWindow) IsAllow(ip string) bool {
	fmt.Println(ip)
	if client, ok := l.clients[ip]; ok {
		allow := client.current < l.cgf.MaxReq
		fmt.Println(client.current)
		if allow {
			l.mu.Lock()
			client.current++
			l.mu.Unlock()
			return allow
		} else {
			return false
		}
	}
	l.mu.Lock()
	c := &counter{
		current: 1,
		resetFunc: func(c *counter, t *time.Timer) {
			<-t.C
			c.current = 0
		},
	}
	l.clients[ip] = c
	c.run(l.cgf.Duration)
	l.mu.Unlock()
	return true
}
