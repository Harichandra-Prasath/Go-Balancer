package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	Url    *url.URL
	Lock   sync.RWMutex
	Proxy  *httputil.ReverseProxy
	Status bool
}

type Pool struct {
	Servers []*Backend
	Current int64
}

func (p *Pool) Addserver(backend *Backend) {
	(p.Servers) = append((p.Servers), backend)
}

func GetPool() *Pool {
	servers := make([]*Backend, 0, 10)

	return &Pool{
		Servers: servers,
		Current: -1,
	}
}

func GetBackend(_url string) *Backend {
	// Assuming when added a backend, it is healthy
	u, _ := url.Parse(_url)
	revproxy := httputil.NewSingleHostReverseProxy(u)
	return &Backend{
		Url:    u,
		Status: true,
		Proxy:  revproxy,
	}

}
