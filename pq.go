package main

import (
	"container/heap"
	"sync"
)

var Lock sync.RWMutex

type Server struct {
	Backend     *Backend
	Connections int
	index       int
}

type Heapq []*Server

func (pq Heapq) Len() int {
	return len(pq)
}

func (pq Heapq) Less(i, j int) bool {
	return pq[i].Connections < pq[j].Connections
}

func (pq Heapq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *Heapq) Push(x any) {
	Lock.Lock()
	defer Lock.Unlock()
	n := len(*pq)
	server := x.(*Server)
	server.index = n
	*pq = append(*pq, server)
}

func (pq *Heapq) Pop() any {
	Lock.Lock()
	defer Lock.Unlock()
	old := *pq
	n := len(old)
	server := old[n-1]
	old[n-1] = nil
	server.index = -1
	*pq = old[0 : n-1]
	return server
}

func (pq *Heapq) update(server *Server, backend *Backend, conn int) {
	Lock.Lock()
	defer Lock.Unlock()
	server.Backend = backend
	server.Connections = conn
	heap.Fix(pq, server.index)
}

func (pq *Heapq) peek() *Server {
	Lock.RLock()
	defer Lock.RUnlock()
	Pq := *pq
	return Pq[0]
}
