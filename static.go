package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
)

func ServeStatic(w http.ResponseWriter, r *http.Request, path *string, mode bool) {

	Conn, _ := r.Context().Value(Key).(net.Conn)
	socket_file, _ := Conn.(*net.TCPConn).File()
	defer socket_file.Close()
	var fn string
	var Root string

	if mode {
		fn = strings.TrimPrefix(*path, "/static/")
		Root = GLOBAL.STATIC_ROOT
	} else {
		fn = strings.TrimPrefix(*path, "/media/")
		Root = GLOBAL.MEDIA_ROOT
	}

	file, err := os.Open(fmt.Sprintf("%s%s", Root, fn))
	info, _ := file.Stat()
	size := info.Size()

	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Requested resource not found\n"))
			Logger.Error(fmt.Sprintf("%s not found in %s", fn, Root))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unknown Error\n"))
			Logger.Error(err.Error())
		}
		return
	}
	go CACHE.add_cache(*path, file, int(size))
	_, err = Conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		panic(err)
	}
	_, err = syscall.Sendfile(int(socket_file.Fd()), int(file.Fd()), nil, int(size))
	if err != nil {
		panic(err)
	}
}
