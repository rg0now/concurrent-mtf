package main

import (
	"fmt"
	"strings"
	"time"
)

// non self-adjusting cache: no cache replacement strategy

type staticCacheNode struct {
	value     Item
	timestamp time.Time
}

type StaticCache struct {
	cache map[Id]staticCacheNode
	store *([]Item)
	size  int
}

func NewStaticCache(store *([]Item)) *StaticCache {
	size := int(CacheSize * float64(len(*store)))
	log("StaticCache: creating cache: size: %d", size)
	c := &StaticCache{store: store, size: size}
	c.cache = make(map[Id]staticCacheNode, size)
	return c
}

func (l *StaticCache) Add(i Item) {
	// log("StaticCache: adding item %d", i.Id())
	// // warm up cache
	// l.Find((*i).Id())
}

func (l *StaticCache) String() string {
	ret := fmt.Sprintf("StaticCache: size %d: ", l.size)
	var ns []string
	for _, n := range l.cache {
		ns = append(ns, fmt.Sprintf("Node: %d (timestamp: %d)",
			n.value.Id(), n.timestamp.UnixMicro()))
	}
	return ret + strings.Join(ns, ", ")
}

func (l *StaticCache) Find(val Item) int {
	// log("StaticCache: find id %d", val)

	// log("before: %s", l.String())

	// cached?
	n, found := l.cache[val.Id()]
	if found {
		n.timestamp = time.Now()
		return 1
	} else {
		// add item to the cache
		i := (*l.store)[val.Id()]

		var id Id = -1
		if len(l.cache) < l.size {
			// add
			id = i.Id()
			// log("StaticCache: new item %d", id)
			newEntry := staticCacheNode{value: i, timestamp: time.Now()}
			l.cache[id] = newEntry
		}

		_, ok := l.cache[i.Id()]
		if ok {
			return 1
		}

		// make cache miss slow
		// time.Sleep(BusyWait)
		res := 0
		for i := 0; i < BusyWait; i++ {
			res += -1
		}

		// log("after: %s", l.String())

		// so that the for loop is not compiled out
		return res
	}
}
