package limiter

import (
	"fmt"
	"sync"
	"time"
)

type FixedWindow struct {
	cfg     Config
	mu      sync.Mutex
	clients map[string]*counter
}

type counter struct {
	current   int
	timeStart time.Time
	timeEnd   time.Time
	resetFunc func(c *counter, t *time.Ticker)
}

func (c *counter) run(dur time.Duration) {
	go c.resetFunc(c, time.NewTicker(dur))
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
		cfg:     cfg,
		mu:      sync.Mutex{},
		clients: make(map[string]*counter),
	}
}

//
func (l *FixedWindow) IsAllow(ip string) bool {
	fmt.Println(ip)
	if client, ok := l.clients[ip]; ok {
		allow := client.current < l.cfg.MaxReq
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
		resetFunc: func(c *counter, t *time.Ticker) {
			for range t.C {
				c.current = 0
			}
		},
	}
	l.clients[ip] = c
	c.run(l.cfg.Duration)
	l.mu.Unlock()
	return true
}
