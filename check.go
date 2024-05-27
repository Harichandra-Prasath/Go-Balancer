package main

import (
	"net"
	"time"
)

// Passive Health check on distinct intervals

func isResponsive(b *Backend) bool {
	conn, err := net.DialTimeout("tcp", b.Url.Host, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return false
}

func (p *Pool) checkAllBackends() {
	// Checks status of all the servers by opening a tcp connection
	for i := range p.Servers {
		backend := p.Servers[i]
		status := isResponsive(backend)
		backend.SetHealth(status)
	}
}

func CheckHealth(p *Pool) {
	t := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-t.C:
			p.checkAllBackends()
		}
	}
}
