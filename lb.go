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

func (p *Pool) GetHealthyBackend() *Backend {
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
