package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var pool *Pool
var Logger *slog.Logger
var Pq Heapq

var ALGO string
var MEDIA_ROOT string
var STATIC_ROOT string

func loadbalancer(w http.ResponseWriter, r *http.Request) {
	Logger.Info(fmt.Sprintf("Client Request at %s", r.URL.Path))
	if strings.HasPrefix(r.URL.Path, "/static/") {
		Logger.Debug("Serving Static Content")
		ServeStatic(w, r, true)
	} else if strings.HasPrefix(r.URL.Path, "/media/") {
		Logger.Debug("Serving Media Content")
		ServeStatic(w, r, false)
	} else {
		var backend *Backend
		switch ALGO {
		case "LC":
			backend = pool.GetHealthyBackendLC()
			defer Release_Connection(backend)
		default:
			backend = pool.GetHealthyBackendRR()
		}
		if backend != nil {

			Logger.Debug(fmt.Sprintf("Proxying the request to %s", backend.Url.Host))
			backend.Proxy.ServeHTTP(w, r)
			return
		} else {
			Logger.Error("No Servers available at the Moment")
			w.Write([]byte("Sorry, No servers available at the Moment\n"))
		}
	}
}

func ConfigLog() {
	Handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	Logger = slog.New(Handler)
}

func main() {

	ConfigLog()
	ALGO = "LC"
	pool = GetPool()
	Logger.Info("Sever Pool Created")

	data, _ := os.ReadFile("backends.txt")
	servers := strings.Split(string(data), "\n")
	Pq = make(Heapq, len(servers))

	for i, server := range servers {
		b := GetBackend(server)
		Pq[i] = &Server{
			Backend:     b,
			Connections: 0,
			index:       i,
		}
		Logger.Debug(fmt.Sprintf("Backend with url %s added to the pool", server))
		pool.Addserver(b)
	}
	Logger.Info("GO-Balancer Started and Serving at 3000")
	go CheckHealth(pool)
	err := http.ListenAndServe(":3000", http.HandlerFunc(loadbalancer))
	if err != nil {
		Logger.Error(err.Error())
		panic(err)
	}

}
