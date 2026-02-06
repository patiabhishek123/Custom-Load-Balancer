package main

import (
	"github.com/patiabhishek123/Custom-Load-Balancer/loadbalancer"
)

func main(){
	// server.RunServer(9)
	loadbalancer.MakeLoadBalancer(5)
}
