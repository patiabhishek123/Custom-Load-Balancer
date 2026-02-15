package circuit

import (
	
	"sync"
	"time"

	"github.com/patiabhishek123/Custom-Load-Balancer/internal/balancer"
)

type Breaker struct {
	failureThreashold int 
	cooldown          time.Duration
	mu                sync.Mutex
}
//
func  NewBreaker(failureThreashold int,cooldown time.Duration) *Breaker{
	return &Breaker{
		failureThreashold: failureThreashold,
		cooldown: cooldown,
	}
}

func (br *Breaker) AllowRequest(b *balancer.Backend) bool {
	state:=b.CircuitState

	if state ==balancer.Open{
		if time.Since(b.LastFailureTime)>br.cooldown{
			b.SetCircuitState(balancer.HalfOpen)
			return true
		}
		return false
	}
	return true

}

func (br *Breaker) RecordFailures(b *balancer.Backend){
	b.IncrementConnections()
	
	if b.FailureCount >=br.failureThreashold{
		b.SetCircuitState(balancer.Open)
		b.SetLastFailureTime(time.Now())
	}
}

func (br *Breaker) RecordSuccess(b *balancer.Backend){
	b.ResetFailure()
	b.SetCircuitState(balancer.Closed)
}