package main

import (
	"fmt"
	"http"
	"log"
	"time"

	//"github.com/patiabhishek123/Custom-Load-Balancer/internal/circuit"
	"github.com/patiabhishek123/Custom-Load-Balancer/internal/balancer"
	"github.com/patiabhishek123/Custom-Load-Balancer/internal/circuit"
	"github.com/patiabhishek123/Custom-Load-Balancer/internal/proxy"
)

// "github.com/patiabhishek123/Custom-Load-Balancer/server"
// "time"

func main() {
	// go server.RunServer(5)

	// Giving servers time to start
	//     time.Sleep(100 * time.Millisecond)

	//     loadbalancer.MakeLoadBalancer(5)

	pool := balancer.NewBackendPool()

	pool.AddBackend(balancer.NewBackend("http://localhost:8081"))
	pool.AddBackend(balancer.NewBackend("http://localhost:8082"))
	pool.AddBackend(balancer.NewBackend("http://localhost:8083"))

	strategy := balancer.NewRoundRobin(pool)
	//or
	//strategy := balancer.NewLeastCount(pool)
	

	// for i := 0; i < 10; i++ {
	// 	b := strategy.NextBackend()
	// 	fmt.Println(b.URL)
	// }

	breaker :=circuit.NewBreaker(3,10*time.Second)
	lb :=proxy.NewLoadBalancer(strategy,breaker)
	log.Println("Load balancer running on :8090")
	log.Fatal(http.ListenAndServe(":8090", lb))
}
