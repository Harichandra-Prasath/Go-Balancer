package main

import (
	"net/http"
	"sync"
)

func make_request(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := http.Get("http://localhost:3000/check")
	if err != nil {
		panic(err)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go make_request(&wg)
	}
	wg.Wait()
}
