package balancer

// pool responsible for storing backends ,thread-safe acess ans routing selection
import (
	"sync"
	"sync/atomic"
)

type BackendPool struct {
	mu       sync.RWMutex // Added a RWMutex as allows Multiple readers and single writer
	backends []*Backend
	counter  uint64
}


func  NewBackendPool() *BackendPool {
	return &BackendPool{
		backends: make([]*Backend , 0),	
	}
}

func (p *BackendPool) AddBackend(b *Backend){
	p.mu.Lock()
	p.backends=append(p.backends,b)
	p.mu.Unlock()
}




func(p *BackendPool) GetHealthyBackends() []*Backend {
	p.mu.RLock()

	healthy:=make([]*Backend, 0)
	for _,b := range p.backends {
		if b.IsHealthy(){
			healthy=append(healthy, b)
		}
	}
	return healthy

}

// func (p *BackendPool) RemoveBackend(b *Backend)


func (p *BackendPool) GetNextBackendRR()*Backend{
	healthy := p.GetHealthyBackends()
	if len(healthy) == 0 {
		return nil
	}
	index := atomic.AddUint64(&p.counter,1)
	return healthy[index%uint64(len(healthy))]
}

func (p *BackendPool) GetAllBackends() []*Backend {
    p.mu.RLock()
    defer p.mu.RUnlock()

    result := make([]*Backend, len(p.backends))
    copy(result, p.backends)
    return result
}
