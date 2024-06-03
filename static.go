package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func ServeStatic(w http.ResponseWriter, path *string, mode bool) {

	var fn string
	var Root string

	if mode {
		fn = strings.TrimPrefix(*path, "/static/")
		Root = GLOBAL.STATIC_ROOT
	} else {
		fn = strings.TrimPrefix(*path, "/media/")
		Root = GLOBAL.MEDIA_ROOT
	}
	content, err := os.ReadFile(fmt.Sprintf("%s%s", Root, fn))
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
	go CACHE.add_cache(*path, content)
	w.Write(content)
}
