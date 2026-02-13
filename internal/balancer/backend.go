package balancer

import (
	// "net/url"
	"sync"
	"time"
)

type Backend struct {
	URL               string
	Healthy           bool
	ActiveConnections int64
	Weight            int
	ResponseTime      time.Duration
	mu                sync.Mutex
}

func NewBackend(url string) *Backend{
	return &Backend{
		URL: url,
		Healthy: true,
		Weight: 1,
	}
}

func (b *Backend) IncrementConnections(){
	b.mu.Lock()
	b.ActiveConnections++
	b.mu.Unlock()
}


func (b *Backend) DecrementConnections(){
	b.mu.Lock()
	if b.ActiveConnections>0{
		b.ActiveConnections--
	}
	b.mu.Unlock()
}


func (b *Backend) IsHealthy() bool{
	b.mu.Lock()
	defer
	b.mu.Unlock()

	return b.Healthy
}

func (b *Backend) SetHealthy(state bool){
	b.mu.Lock()
	if b.Healthy!=state{
		b.Healthy=state
	}
	b.mu.Unlock()
}





