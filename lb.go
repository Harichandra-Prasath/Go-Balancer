package main

import (
	"math"
	"sync/atomic"
)

func (b *Backend) SetHealth(status bool) {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	b.Status = status

}

func (b *Backend) getHealth() bool {
	b.Lock.RLock()
	defer b.Lock.RUnlock()
	return b.Status
}

func (p *Pool) GetCurrent() int {
	// Returns the next index which contains the backend underneath
	// Usage of atomic is to avoid if there are multiple connections,access this function at the same time
	return int(atomic.AddInt64(&p.Current, int64(1)) % int64(len(p.Servers)))
}

func (p *Pool) GetHealthyBackendRR() *Backend {
	// This functions implies the balancing algorithm (Round Robin) for backends and return them if and only if they are healthy
	start := p.GetCurrent()
	total := len(p.Servers)
	width := start + total

	for i := start; i < width; i++ {
		current := i % len(p.Servers) // Padding
		if p.Servers[current].getHealth() {

			if i != start {
				atomic.AddInt64(&p.Current, int64(i))
			}
			return p.Servers[current]
		}
	}
	return nil // No servers or none of them are healthy
}

func (p *Pool) GetHealthyBackendLC() *Backend {
	// This functions implies the balancing algorithm (Least Connections) for backends and return them if and only if they are healthy
	for len(Pq) > 0 {

		server := Pq.peek()
		backend := server.Backend
		if backend.getHealth() {
			Pq.update(server, backend, server.Connections+1)
			return backend
		} else {
			return nil
		}
	}

	return nil
}

func Release_Connection(b *Backend) {
	if b != nil {
		s := Get_server(b)
		Pq.update(s, b, s.Connections-1)
	}

}

func Get_server(b *Backend) *Server {
	for _, s := range Pq {
		if s.Backend == b {
			return s
		}
	}
	return nil // no chance
}

func Update_health(b *Backend, health bool) {
	if b.getHealth() == health {
		return
	} else {
		s := Get_server(b)
		if b.getHealth() {
			Pq.update(s, b, math.MaxInt)
		} else {
			Pq.update(s, b, 0)
		}
	}

}
