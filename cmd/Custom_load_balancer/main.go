package main

import (
    "github.com/patiabhishek123/Custom-Load-Balancer/loadbalancer"
    "github.com/patiabhishek123/Custom-Load-Balancer/server"
    "time"
)

func main(){
    go server.RunServer(5)
    
    // Giving servers time to start
    time.Sleep(100 * time.Millisecond)
    
    loadbalancer.MakeLoadBalancer(5)
}