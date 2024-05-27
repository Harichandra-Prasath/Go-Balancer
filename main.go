package main

import (
	"os"
	"strings"
)

func main() {
	pool := GetPool()

	data, _ := os.ReadFile("backends.txt")
	servers := strings.Split(string(data), "\n")

	for _, server := range servers {
		b := GetBackend(server)
		pool.Addserver(b)
	}
}
