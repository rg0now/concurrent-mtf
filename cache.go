package main

import (
        "fmt"        
        "strings"
        "time"
)

const BusyWait = 100 * time.Microsecond
const CacheSize float64 = 0.2

type lruCacheNode struct {
        value *Item
        timestamp time.Time
}

type LruCache struct {
        cache map[Id]lruCacheNode
        store *([]Item)
        size int
}

func NewLruCache(store *([]Item)) *LruCache {
        size := int(CacheSize * float64(len(*store)))
        log("LruCache: creating cache: size: %d", size)
        c := &LruCache{store: store, size: size}
        c.cache = make(map[Id]lruCacheNode, size)
        return c
}

func (l *LruCache) Add(i *Item) {
        log("LruCache: adding item %d", (*i).Id())
        // // warm up cache
        // l.Find((*i).Id())
}

func (l *LruCache) String() string {
        ret := fmt.Sprintf("LruCache: size %d: ", l.size)
        var ns []string
        for _, n := range l.cache {
		ns = append(ns, fmt.Sprintf("Node: %d (timestamp: %d)",
                        (*n.value).Id(), n.timestamp.UnixMicro()))
        }
        return ret + strings.Join(ns, ", ")
}

func (l *LruCache) Find(val Id) int {
        log("LruCache: find id %d", val)

        // log("before: %s", l.String())

        // cached?
        n, found := l.cache[val]
        if found {
                n.timestamp = time.Now()
                return 1
        } else {
                // add item to the cache
                i := (*l.store)[val]
                
                lru := time.Now()
                var id Id = -1
                if len(l.cache) < l.size {
                        // add
                        id = i.Id()
                        log("LruCache: new item %d", id)
                } else {
                        // replace
                        for k, e := range l.cache {
                                if e.timestamp.Before(lru) {
                                        id = k
                                        lru = e.timestamp
                                }
                        }

                        log("LruCache: replace item: %d -> %d", id, i.Id())
                        delete(l.cache, id)
                        id = i.Id()
                }
                newEntry := lruCacheNode{value: &i, timestamp: time.Now() }
                l.cache[id] = newEntry

                // make cache miss slow
                time.Sleep(BusyWait)

                // log("after: %s", l.String())

                return -1
        }
}
