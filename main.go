package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	flag "github.com/spf13/pflag"
)

const BufferSize int = 256_000
const CacheBusyWait = 100_000
const WeightedTreeBusyWait = 100_000
const DefaultCacheHitRate float64 = 0.1
const BloomFalsePositiveRate float64 = 0.001
const BloomFillRatio float64 = 0.5

var CacheHitRate float64

var verbose bool = false
var seed int64

func log(format string, v ...any) {
	if verbose {
		fmt.Printf(format+"\n", v...)
	}
}

type Id = int
type Channel chan Item

type Item interface {
	Less(Item) bool
	Id() Id
	Match(i Id) bool
}

type Source interface {
	Generate() (Item, error)
}

type LoadBalancer interface {
	Generate() error
	Assign(i Item)
}

type DataStructure interface {
	Add(p Item)
	Find(i Item) int
	String() string
}

type Thread struct {
	id int
	rx Channel
	d  DataStructure
}

var ThreadStore []*Thread
var LBStore []LoadBalancer
var isBloom bool

func main() {
	// general flags
	os.Args[0] = "cmtf"
	var n = flag.IntP("req_num", "n", 10, "Number of requests")
	var m = flag.IntP("item_num", "m", 5, "Number of items")
	var k = flag.IntP("thread_num", "k", 1, "Number of threads")
	var lk = flag.IntP("lb_thread_num", "t", 1, "Number of LB threads")
	var ds = flag.StringP("data-structure", "d", "mtf",
		"Data-structure: cache|nullcache|mtf|linkedlist|splay|btree|wsplay|wbtree|bloom|sabloom|wbloom|wsabloom")
	var src = flag.StringP("source", "s", "uniform", "Source: uniform|poisson|zipf:a")
	var sp = flag.StringP("load-balancer", "l", "modulo", "Load-balaner: modulo|split|roundrobin")
	var v = flag.BoolP("verbose", "v", false, "Verbose logging, identical to <-l all:DEBUG>")
	flag.Int64Var(&seed, "seed", 1, "Seed")
	flag.Float64Var(&CacheHitRate, "cache-hit-rate", DefaultCacheHitRate, "Target cache hit rate")
	flag.Parse()
	verbose = *v
	isBloom = strings.HasSuffix(strings.ToLower(*ds), "bloom")

	log("Creating & shuffling item list")
	is := make([]Item, *m)
	for j := 0; j < *m; j++ {
		is[j] = IntegerItem{j}
	}

	log("Creating comm channels")
	cs := make([]Channel, *k)
	for j := 0; j < *k; j++ {
		cs[j] = make(Channel, BufferSize)
	}

	log("Creating LB(s) and source(s)")
	LBStore = make([]LoadBalancer, *lk)
	for i := 0; i < *lk; i++ {
		var s Source
		switch {
		case *src == "uniform":
			s = NewUniformSource((*m)/(*lk), *n, &is)
		case *src == "poisson":
			s = NewPoissonSource((*m)/(*lk), *n, float64(*m)/2.0, &is)
		case strings.HasPrefix(*src, "zipf"):
			s = NewZipfSource(*src, (*m)/(*lk), *n, &is)
		default:
			panic("Unknown source type: " + *src)
		}

		var lb LoadBalancer
		switch strings.ToLower(*sp) {
		case "modulo":
			lb = NewModuloLB(*k, s, cs)
		case "split":
			lb = NewSplitLB(*k, *m, s, cs)
		case "roundrobin":
			lb = NewRoundRobinLB(*k, s, cs)
		default:
			panic("Unknown load-balancer: " + *sp)
		}
		LBStore[i] = lb
	}

	log("Initializing threads")
	ThreadStore = make([]*Thread, *k)
	wg := new(sync.WaitGroup)
	wg.Add(*k)
	for j := 0; j < *k; j++ {
		// create data-structure
		var d DataStructure
		switch strings.ToLower(*ds) {
		case "cache":
			d = NewLruCache(*m, true)
		case "nullcache", "statcache":
			d = NewLruCache(*m, false)
		case "mtf":
			d = NewMtf()
		case "linkedlist":
			d = NewLinkedList()
		case "splay":
			d = NewSplayTree()
		case "btree":
			d = NewBTree()
		case "wsplay":
			d = NewWeightedSplayTree()
		case "wbtree":
			d = NewWeightedBTree()
		case "bloom":
			d = NewBloomFilter(*m, BloomFalsePositiveRate)
		case "sabloom":
			d = NewSelfAdjustingBloomFilter(*m, BloomFalsePositiveRate)
		case "wbloom":
			d = NewWeightedBloomFilter(*m, BloomFalsePositiveRate)
		case "wsabloom":
			d = NewWeightedSelfAdjustingBloomFilter(*m, BloomFalsePositiveRate)
		default:
			panic("Unknown data structure: " + *ds)
		}

		// wrap data-structure with a thread
		t := &Thread{id: j, rx: cs[j], d: d}
		ThreadStore[j] = t

		// spawn thread for initialization
		go func(t *Thread) {
			for i := 0; i < *m; i++ {
				t.d.Add(is[i])
				// if verbose {
				// 	log("Thread %d: added: %d, state: %s",
				// 		t.id, is[i].Id(), t.d.String())
				// }
			}
			wg.Done()
		}(t)
	}
	wg.Wait()

	log("Initialization done, starting workers")
	wg.Add(*k)
	for j := 0; j < *k; j++ {
		// spawn worker thread
		go func(t *Thread) {
			l := t.d
			found := 0
			tt0 := time.Now()

			for r := range t.rx {
				// if verbose {
				// 	log("Thread %d: %s", t.id, l.String())
				// }
				if verbose && found == 0 {
					tt0 = time.Now() // reset thread timer: this is the real point where thread starts
				}

				i := l.Find(r)
				if i >= 0 {
					found++
				}
				if verbose {
					log("Thread %d: lookup: %d, found: %t, state: %s",
						t.id, r.Id(), i > 0, l.String())
				}
			}
			wg.Done()

			// fmt.Printf("Thread %d exited, found: %d\n", t.id, found)

			tt1 := time.Now()
			dt := tt1.Sub(tt0)

			fmt.Printf("Thread %d exited, found: %d, running time: %s\n", t.id,
				found, dt)

			log("Thread %d exited, found: %d, running time: %s", t.id, found, dt)

			return
		}(ThreadStore[j])
	}

	log("Starting main loop")
	t0 := time.Now()
	wglb := new(sync.WaitGroup)
	wglb.Add(*lk)
	for i := 0; i < *lk; i++ {
		go func(lb LoadBalancer) {
			for {
				err := lb.Generate()
				if err != nil {
					break
				}
			}
			wglb.Done()
		}(LBStore[i])
	}

	// close all channels
	wglb.Wait()
	for j := 0; j < *k; j++ {
		close(cs[j])
	}

	// wait for all threads to finish
	wg.Wait()

	t1 := time.Now()
	d := t1.Sub(t0)

	if isBloom {
		totalHashNum := 0
		for j := 0; j < *k; j++ {
			l := ThreadStore[j].d
			// fmt.Printf("Final state: %s\n", l.String())
			if isBloom {
				switch b := l.(type) {
				case *BloomFilter:
					totalHashNum += b.hashCounter
				case *WeightedBloomFilter:
					totalHashNum += b.hashCounter
				case *WeightedSelfAdjustingBloomFilter:
					totalHashNum += b.hashCounter
				}
			}
		}
		fmt.Printf("%d\t%d\t%d\t%v\t%f\t%d\n", *k, *m, *n, d, float64(*n)/d.Seconds(), totalHashNum)
	} else {

		fmt.Printf("%d\t%d\t%d\t%v\t%f\n", *k, *m, *n, d, float64(*n)/d.Seconds())
	}
}
