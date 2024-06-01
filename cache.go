package main

import "sync"

// Caching for static and media files
type Cache map[string][]byte

var CLock sync.RWMutex

func (c *Cache) add_cache(path string, content []byte) {
	CLock.Lock()
	defer CLock.Unlock()
	C := *c
	C[path] = content
}

func (c *Cache) get_cache(path string) (bool, []byte) {
	CLock.RLock()
	defer CLock.RUnlock()
	C := *c
	content, isCached := C[path]
	return isCached, content
}
