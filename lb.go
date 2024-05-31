package main

import (
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

func (p *Pool) Schedule() *Backend {
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
			atomic.AddInt64(&p.Active_Connections, int64(1))
			return p.Servers[current]
		}
	}
	return nil // No servers or none of them are healthy
}

func (P *Pool) Release_Connection(b *Backend) {
	Logger.Debug("Request Completed: Removing Connection Registry")
	atomic.AddInt64(&P.Active_Connections, int64(-1))
}

func (Pq *Heapq) Schedule() *Backend {
	// This functions implies the balancing algorithm (Least Connections) for backends and return them if and only if they are healthy
	for len(*Pq) > 0 {

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

func (Pq *Heapq) Release_Connection(b *Backend) {
	Logger.Debug("Request Completed: Removing Connection Registry")
	s := Pq.Get_server(b)
	Pq.update(s, b, s.Connections-1)
}

func (Pq *Heapq) Get_server(b *Backend) *Server {
	for _, s := range *Pq {
		if s.Backend == b {
			return s
		}
	}
	return nil // no chance
}
