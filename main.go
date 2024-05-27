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

func loadbalancer(w http.ResponseWriter, r *http.Request) {
	Logger.Info(fmt.Sprintf("Client Request at %s", r.URL.Path))
	if strings.HasPrefix(r.URL.Path, "/static/") {
		Logger.Debug("Serving Static Content")
	} else if strings.HasPrefix(r.URL.Path, "/media/") {
		Logger.Debug("Serving Media Content")
	} else {
		backend := pool.GetHealthyBackend()
		if backend != nil {
			Logger.Debug(fmt.Sprintf("Proxying the request to %s", backend.Url.Host))
			backend.Proxy.ServeHTTP(w, r)
			return
		} else {
			Logger.Error("No Servers available at the Moment")
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

	pool = GetPool()
	Logger.Info("Sever Pool Created")

	data, _ := os.ReadFile("backends.txt")
	servers := strings.Split(string(data), "\n")

	for _, server := range servers {
		b := GetBackend(server)
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
