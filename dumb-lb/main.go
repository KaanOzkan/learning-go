package main

import (
    "fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

// Core object
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy // Takes incoming requests and sends to another server
}

type ServerPool struct {
	backends []*Backend
	current  uint64
}

var serverPool ServerPool

func main() {
	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(u)
    server:= http.Server{
        Addr: fmt.Sprintf(":%d", 3030),
        Handler: http.HandlerFunc(reverseProxy.ServeHTTP),
    }
}

func loadBalance(w http.ResponseWriter, r *http.Request) {
    peer := serverPool.GetNextPeer()
    if peer != nil {
        peer.ReverseProxy.ServeHTTP(w, r)
        return
    }
    http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}

func (s *ServerPool) GetNextPeer() *Backend {
    next := s.NextIndex()
    full_cycle := len(s.backends) + next
    for i := next; i < full_cycle; i++ {
        idx := i % len(s.backends)
        if s.backends[idx].IsAlive() {
            if i != next { // not the original backend
                atomic.StoreUint64(&s.current, uint64(idx)) // mark
            }
            return s.backends[idx]
        }
    }
    return nil
}

func (s *ServerPool) NextIndex() int {
    // increase value by one, return that index.
    return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

func (b *Backend) SetAlive(alive bool) {
    b.mux.Lock()
    b.Alive = alive
    b.mux.Unlock()
}

func (b *Backend) IsAlive() (alive bool) {
    b.mux.RLock()
    alive = b.Alive
    b.mux.RUnlock()
    return
}
