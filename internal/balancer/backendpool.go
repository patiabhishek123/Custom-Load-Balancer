package balancer

import (
	"sync"
)

type BackendPool struct {
	mu       sync.RWMutex // Added a RWMutex as allows Multiple readers and single writer
	backends []int
	counter  uint64
}

func (p *BackendPool) AddBackend(b *Backend)

func (p *BackendPool) RemoveBackend(b *Backend)

func(p *BackendPool) GetHealthyBackends() []*Backend

func GetNExtBackendRR()