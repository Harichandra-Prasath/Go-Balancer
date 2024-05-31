package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

var MANAGER Manager
var Logger *slog.Logger

func loadbalancer(w http.ResponseWriter, r *http.Request) {
	Logger.Info(fmt.Sprintf("Client Request at %s", r.URL.Path))
	if strings.HasPrefix(r.URL.Path, "/static/") {
		Logger.Debug("Serving Static Content")
		ServeStatic(w, r, true)
	} else if strings.HasPrefix(r.URL.Path, "/media/") {
		Logger.Debug("Serving Media Content")
		ServeStatic(w, r, false)
	} else {
		backend := MANAGER.Schedule()
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

func main() {
	err := InitialiseSystem()
	if err != nil {
		panic(err)
	}
	Logger.Info(fmt.Sprintf("GO-Balancer Started and Serving at %d", GLOBAL.Port))
	go CheckHealth(MANAGER)
	err = http.ListenAndServe(fmt.Sprintf(":%d", GLOBAL.Port), http.HandlerFunc(loadbalancer))
	if err != nil {
		Logger.Error(err.Error())
		panic(err)
	}

}
