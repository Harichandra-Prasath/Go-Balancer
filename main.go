package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

var MANAGER Manager
var CACHE Cache
var Logger *slog.Logger

type ConnCtx struct {
	key string
}

var Key = &ConnCtx{
	key: "underlying-conn",
}

func setContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, Key, c)
}

func loadbalancer(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	Logger.Info(fmt.Sprintf("Client Request at %s", path))

	Is_Static := strings.HasPrefix(path, "/static/")
	Is_Media := strings.HasPrefix(path, "/media/")

	if Is_Static || Is_Media {
		cached, content := CACHE.get_cache(path)
		if cached {
			Logger.Debug("Cache Available. Reading from Cache for the request")
			w.Write(content)
			return
		} else {
			if Is_Static {
				ServeStatic(w, r, &path, true)
			} else {
				ServeStatic(w, r, &path, false)
			}
		}
	} else {
		backend := MANAGER.Schedule(GLOBAL.ALGO)
		if backend != nil {
			defer MANAGER.Release_Connection(backend)
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
	server := http.Server{
		Addr:        fmt.Sprintf(":%d", GLOBAL.Port),
		Handler:     http.HandlerFunc(loadbalancer),
		ConnContext: setContext,
	}
	err = server.ListenAndServe()
	if err != nil {
		Logger.Error(err.Error())
		panic(err)
	}

}
