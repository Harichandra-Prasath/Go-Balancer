package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Manager interface {
	Addserver(*Backend)
	Schedule() *Backend
	checkAllBackends()
	UpdateHealth(*Backend, bool)
	Release_Connection(*Backend)
}

type Backend struct {
	Url    *url.URL
	Lock   sync.RWMutex
	Proxy  *httputil.ReverseProxy
	Status bool
}

type Pool struct {
	Servers            []*Backend
	Current            int64
	Active_Connections int64
}

func (p *Pool) Addserver(backend *Backend) {
	(p.Servers) = append((p.Servers), backend)
}

func GetPool(n int) *Pool {
	servers := make([]*Backend, n)

	return &Pool{
		Servers: servers,
		Current: -1,
	}
}

func GetQueue(n int) *Heapq {
	Queue := make(Heapq, n)
	return &Queue
}

func (pq *Heapq) Addserver(b *Backend) {
	len := len(*pq)
	*pq = append(*pq, &Server{
		Backend:     b,
		Connections: 0,
		index:       len,
	})
}

func GetBackend(_url string) *Backend {
	// Assuming when added a backend, it is healthy
	u, _ := url.Parse(_url)
	revproxy := httputil.NewSingleHostReverseProxy(u)
	return &Backend{
		Url:    u,
		Status: false,
		Proxy:  revproxy,
	}

}
