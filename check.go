package main

import (
	"fmt"
	"math"
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
	MANAGER.UpdateHealth(b, healthy)
	return healthy

}

func (p *Pool) checkAllBackends() {
	// Checks status of all the servers by opening a tcp connection
	for i := range p.Servers {
		backend := p.Servers[i]
		isResponsive(backend)
	}
}

func (Pq *Heapq) checkAllBackends() {
	// Checks status of all the servers by opening a tcp connection
	for _, server := range *Pq {
		isResponsive(server.Backend)
	}
}

func CheckHealth(MANAGER Manager) {
	t := time.NewTicker(time.Duration(GLOBAL.Poll_Period) * time.Second)
	for {
		select {
		case <-t.C:
			Logger.Debug("Health Check Started")
			MANAGER.checkAllBackends()
			Logger.Debug("Health Check Completed")
		}
	}
}

func (Pq *Heapq) UpdateHealth(b *Backend, health bool) {
	if b.getHealth() == health {
		return
	} else {
		s := Pq.Get_server(b)
		if b.getHealth() {
			Pq.update(s, b, math.MaxInt)
		} else {
			Pq.update(s, b, 0)
		}
	}
	b.SetHealth(health)
}

func (pool *Pool) UpdateHealth(b *Backend, health bool) {
	b.SetHealth(health)
}
