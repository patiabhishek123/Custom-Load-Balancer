package balancer

import (
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL               url.URL
	Healthy           bool
	ActiveConnections int64
	Weight            int
	ResponseTime      time.Duration
	mu                sync.Mutex
}

func (b *Backend) IncrementConnections()
func (b *Backend) DecrementConnections()
func (b *Backend) IsHealthy() bool
func (b *Backend) SetHealthy()

