package server

import (
	"fmt"
	"log"
)


type ServerList struct{
	Ports[] int
}

func (s *ServerList) Populate(amount int) {
	if amount>=10{
		log.Fatal("Cannot use more than 10 servers")
	}
	for i := 0; i < amount; i++ {
		s.Ports = append(s.Ports, 8000+i)
	}
}

func (s *ServerList) Pop() int{
	port:=s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}
func Test(){
	myServerList := ServerList{}
	myServerList.Populate(5)
	x:=len(myServerList.Ports)
	for i := 0; i < x; i++ {
		fmt.Println(myServerList.Pop())
	}
}