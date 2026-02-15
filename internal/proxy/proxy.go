package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"fmt"
	"github.com/patiabhishek123/Custom-Load-Balancer/internal/balancer"
	"github.com/patiabhishek123/Custom-Load-Balancer/internal/circuit"
)


type LoadBalancer struct {
	strategy balancer.Strategy
	breaker *circuit.Breaker
}

func NewLoadBalancer(strategy balancer.Strategy,breaker *circuit.Breaker)*LoadBalancer{
	return &LoadBalancer{
		strategy: strategy,
		breaker: breaker,
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request){
	backend := lb.strategy.NextBackend()
	fmt.Println("Forwarding to:", backend.URL)

	if backend == nil{
		http.Error(w,"No healthy backends",http.StatusServiceUnavailable)
		return 
	}

	if !lb.breaker.AllowRequest(backend){
		http.Error(w,"Backend temporarily unavailable",http.StatusServiceUnavailable)
		return
	}

	target,_:= url.Parse(backend.URL)
	proxy:=httputil.NewSingleHostReverseProxy(target)

	proxy.ErrorHandler=func(w http.ResponseWriter, r *http.Request, err error) {
		lb.breaker.RecordFailures(backend)
		http.Error(w,"Upstream error",http.StatusBadGateway)
	}

	backend.IncrementConnections()
	defer backend.DecrementConnections()

	proxy.ServeHTTP(w,r)
	lb.breaker.RecordSuccess(backend)
}