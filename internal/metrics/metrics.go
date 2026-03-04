package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/patiabhishek123/Custom-Load-Balancer/internal/balancer"
)

type PoolMetrics struct {
	TotalRequests     uint64 `json:"total_requests"`
	FailedRequests    uint64 `json:"failed_requests"`
	ActiveConnections uint64 `json:"active_connections"`
}

var (
	mu    sync.RWMutex
	pools = make(map[*balancer.BackendPool]*PoolMetrics)
)

func RegisterPool(p *balancer.BackendPool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := pools[p]; !ok {
		pools[p] = &PoolMetrics{}
	}
}

func getMetrics(p *balancer.BackendPool) *PoolMetrics {
	mu.RLock()
	m := pools[p]
	mu.RUnlock()
	return m
}

func IncTotalRequests(p *balancer.BackendPool) {
	m := getMetrics(p)
	if m == nil {
		return
	}
	atomic.AddUint64(&m.TotalRequests, 1)
}

func IncFailedRequests(p *balancer.BackendPool) {
	m := getMetrics(p)
	if m == nil {
		return
	}
	atomic.AddUint64(&m.FailedRequests, 1)
}

func IncActiveConnections(p *balancer.BackendPool) {
	m := getMetrics(p)
	if m == nil {
		return
	}
	atomic.AddUint64(&m.ActiveConnections, 1)
}

func DecActiveConnections(p *balancer.BackendPool) {
	m := getMetrics(p)
	if m == nil {
		return
	}
	// avoid underflow
	for {
		old := atomic.LoadUint64(&m.ActiveConnections)
		if old == 0 {
			return
		}
		if atomic.CompareAndSwapUint64(&m.ActiveConnections, old, old-1) {
			return
		}
	}
}

// Handler returns JSON with metrics for all registered pools.
func Handler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	out := make(map[string]*PoolMetrics)
	for p, m := range pools {
		// use pool pointer address as the key
		key := fmt.Sprintf("%p", p)
		// copy snapshot
		out[key] = &PoolMetrics{
			TotalRequests:     atomic.LoadUint64(&m.TotalRequests),
			FailedRequests:    atomic.LoadUint64(&m.FailedRequests),
			ActiveConnections: atomic.LoadUint64(&m.ActiveConnections),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
