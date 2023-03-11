package metastore

import (
	"sync"
)

type InMemoryBackend struct {
	indexMap sync.Map
}

var imb *InMemoryBackend
var imbOnce sync.Once

func GetInMemoryBacked() *InMemoryBackend {
	imbOnce.Do(func() {
		imb = &InMemoryBackend{}
	})
	return imb
}

func (imb *InMemoryBackend) Get(key string) (any, bool) {
	return imb.indexMap.Load(key)
}

func (imb *InMemoryBackend) Set(key string, value bool) (any, bool) {
	imb.indexMap.Store(key, value)
	return nil, true
}
