package server

import (
	"fmt"
	"net/http"
	"sync"
)

type ServerList struct {
	mu    sync.Mutex // Added a Mutex for thread safety
	Ports []int
}

func (s *ServerList) Populate(amount int) {
	for i := 0; i < amount; i++ {
		s.Ports = append(s.Ports, 8000+i)
	}
}

func (s *ServerList) Pop() int {
	s.mu.Lock()         // Lock before touching the slice
	defer s.mu.Unlock() // Unlock when done
	
	if len(s.Ports) == 0 {
		return 0
	}
	port := s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}

func RunServer(amount int) {
	myServerList := ServerList{}
	myServerList.Populate(amount)

	var wg sync.WaitGroup
	wg.Add(amount)

	for i := 0; i < amount; i++ {
		// Pass wg as a pointer (&wg)
		go MakeServer(&myServerList, &wg)
	}
	
	// This keeps the main process alive!
	wg.Wait() 
}

func MakeServer(s1 *ServerList, wg *sync.WaitGroup) {
	defer wg.Done() // Signal when the server stops
	port := s1.Pop()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from port %d", port)
	})

	// Corrected formatting from "&d" to ":%d"
	addr := fmt.Sprintf(":%d", port) 
	fmt.Printf("Starting server on %s\n", addr)
	
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server on %s failed: %s\n", addr, err)
	}
}