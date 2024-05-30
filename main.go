package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var manager Manager
var Logger *slog.Logger

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
		backend := manager.Schedule()
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
	Logger.Info("Sever Pool Created")

	data, _ := os.ReadFile("backends.txt")
	servers := strings.Split(string(data), "\n")
	fmt.Println(servers)
	Logger.Info("GO-Balancer Started and Serving at 3000")
	go CheckHealth(manager)
	err := http.ListenAndServe(":3000", http.HandlerFunc(loadbalancer))
	if err != nil {
		Logger.Error(err.Error())
		panic(err)
	}

}
