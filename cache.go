package main

import (
        "fmt"        
        "strings"

        lru "github.com/hashicorp/golang-lru"
)

type LruCache struct {
        selfadjusting bool
        m, size int
        cache *lru.Cache
}

// static cache never evicts
func NewLruCache(m int, selfadjusting bool) *LruCache {
        size := int(CacheSize * float64(m))
        log("LruCache: creating cache: size: %d", size)
        cache, err := lru.New(size)
        if err != nil {
                panic(err)
        }
        c := &LruCache{m: m, cache: cache, size: size, selfadjusting: selfadjusting}
        return c
}

func (l *LruCache) Add(i Item) {
        log("LruCache: adding item %d", i.Id())
        // // warm up cache
        // _ = l.cache.Add(i, i)
}

func (l *LruCache) String() string {
        ret := fmt.Sprintf("LruCache: size %d: ", l.size)
        var ns []string
        for  n := 0; n < l.m; n += 1 {
                if l.cache.Contains(IntegerItem{id: n}) {
                        ns = append(ns, fmt.Sprintf("%d", n))
                }
        }
        return ret + strings.Join(ns, ", ")
}

func (l *LruCache) Find(val Item) int {

        // cached?
        var found bool
        if l.selfadjusting {
                log("LruCache: find id %d", val)
                _, ok := l.cache.Get(val)
                if !ok {
                        log("LruCache: Add item %d", val.Id())
                        l.cache.Add(val, val)
                }
                found = ok
        } else {
                log("StatCache: find id %d", val)
                _, found = l.cache.Peek(val)
                if l.cache.Len() <  l.size {
                        log("StatCache: Add item %d", val.Id())
                        l.cache.Add(val, val)
                }
        }
        if found {
                log("LruCache: Found %d", val.Id())
                return 1
        } else {
                log("LruCache: Cache miss for %d, busy wait", val.Id())
                res := 0
                for i := 0; i < BusyWait; i++ {
                        res += -1
                }

                log("LruCache: Cache miss for %d busy wait over", val.Id())
                // so that the for loop is not compiled out
                return res
        }
}
