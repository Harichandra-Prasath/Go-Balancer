package main

import (
	"net/http"
	"os"
	"strings"
)

var pool *Pool

func loadbalancer(w http.ResponseWriter, r *http.Request) {
	backend := pool.GetHealthyBackend()
	if backend != nil {
		backend.Proxy.ServeHTTP(w, r)
		return
	}
}

func main() {
	pool = GetPool()

	data, _ := os.ReadFile("backends.txt")
	servers := strings.Split(string(data), "\n")

	for _, server := range servers {
		b := GetBackend(server)
		pool.Addserver(b)
	}

	err := http.ListenAndServe(":3000", http.HandlerFunc(loadbalancer))
	if err != nil {
		panic(err)
	}

	go CheckHealth(pool)
}
