package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

var (
	baseURL = "http://localhost:808"
)

type RoundRobin struct {
	RevProxy httputil.ReverseProxy
}

type Endpoints struct {
	List []*url.URL
}

func (e *Endpoints) Shuffle() {
	temp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, temp)
}

func MakeLoadBalancer(amount int) {
	// Instantiate Objects
	 lb:= RoundRobin{}
	 ep:= Endpoints{}

	// Server + Router
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	// Creating the endpoints
	for i := 0; i < amount; i++ {
		ep.List = append(ep.List, createEndpoint(baseURL, i))
	}

	// Handler Functions
	router.HandleFunc("/loadbalancer", makeRequest(&lb, &ep))

	// Listen and Server
	log.Fatal(server.ListenAndServe())
}

func makeRequest(lb *RoundRobin, ep *Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for !testServer(ep.List[0].String()) {
			ep.Shuffle()
		}
		lb.RevProxy = *httputil.NewSingleHostReverseProxy(ep.List[0])
		ep.Shuffle()
		lb.RevProxy.ServeHTTP(w, r)
	}
}

func createEndpoint(endpoint string, idx int) *url.URL {
	link := endpoint + strconv.Itoa(idx)
	url, _ := url.Parse(link)
	return url
}

func testServer(endpoint string) bool {
	resp, err := http.Get(endpoint)
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}