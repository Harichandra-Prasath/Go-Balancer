package main

import (
	"os"
	"sync"
)

// Caching for static and media files
type Cache map[string][]byte

var CLock sync.RWMutex

func (c *Cache) add_cache(path string, file *os.File, size int) {
	defer file.Close()
	CLock.Lock()
	defer CLock.Unlock()
	C := *c
	buff := make([]byte, size)
	n, _ := file.Read(buff)
	C[path] = buff[:n]
	Logger.Debug("Writing into Cache")
}

func (c *Cache) get_cache(path string) (bool, []byte) {
	CLock.RLock()
	defer CLock.RUnlock()
	C := *c
	content, isCached := C[path]
	return isCached, content
}
