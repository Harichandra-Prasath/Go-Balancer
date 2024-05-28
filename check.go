package main

import (
	"fmt"
	"net"
	"time"
)

// Passive Health check on distinct intervals

func isResponsive(b *Backend) bool {
	conn, err := net.DialTimeout("tcp", b.Url.Host, 3*time.Second)
	healthy := false
	if err != nil {
		Logger.Warn(fmt.Sprintf("Server with url %s is Unresponsive", b.Url.Host))
	} else {
		conn.Close()
		healthy = true
	}
	if ALGO == "LC" {
		Update_health(b, healthy)
	}
	return healthy

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
			Logger.Debug("Health Check Started")
			p.checkAllBackends()
			Logger.Debug("Health Check Completed")
		}
	}
}
