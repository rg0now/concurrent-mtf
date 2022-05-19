package main

import (
        "fmt"
        "time"
        "os"
        "strings"
        "sync"

        flag "github.com/spf13/pflag"
)


const BufferSize int = 30
//const BufferSize int = 1

var verbose bool = false

func log(format string, v ...any) {
        if verbose {
                fmt.Printf(format + "\n", v...)
        }
}

type Id = int
type Channel chan Id

type Item interface {
        Id() Id
        Match(i Id) bool
}

type Source interface {
        Generate() (Id, error)
}

type LoadBalancer interface {
        Assign(i Id)
}

type DataStructure interface {
        Add(p *Item)
        Find(i Id) int
        String() string
}

type Thread struct {
        rx Channel
        d DataStructure
}

func main() {
        // general flags
        os.Args[0] = "cmtf"
        var n       = flag.IntP("req_num", "n", 10, "Number of requests.")
        var m       = flag.IntP("item_num", "m", 5, "Number of items.")
        var k       = flag.IntP("thread_num", "k", 1, "Number of threads.")
        var ds      = flag.StringP("data-structure", "d", "mtf", "Data-structure: mtf|cache|splay (default: mtf).")
        var v       = flag.BoolP("verbose", "v", false, "Verbose logging, identical to <-l all:DEBUG>.")
        flag.Parse()
        verbose = *v

        // source
        s := UniformSource {
                n: *n,
                m: *m,
                i: 0,
        }

        // items
        is := make([]Item, *m)
        for j := 0; j < *m; j++ {
                is[j] = IntegerItem{j}
        }
        
        // comm channels
        cs := make([]Channel, *k)
        for j := 0; j < *k; j++ {
                cs[j] = make(Channel, BufferSize)
        }

        // loadbalancer
        lb := NewModuloLB(*k, cs)

        // init threads
        wg := new(sync.WaitGroup)
        wg.Add(*k)
        for j := 0; j < *k; j++ {
                // create data-structure
                var d DataStructure
                switch strings.ToLower(*ds) {
                case "mtf": d   = NewMtf()
                case "cache": d = NewLruCache(&is)
                default: panic("Unknown data structure: " + *ds)
                }
                
                for i := 0; i < *m; i++ {
                        d.Add(&is[i])
                }
                // wrap data-structure with a thread
                log("Thread: adding new thread: %d", j)
                t := &Thread{rx: cs[j], d: d}

                // spawn thread
                go func() {
                        l := t.d
                        found := 0
                        for r := range t.rx {
                                if verbose {
                                        log("Thread %d: %s", j, l.String())
                                }

                                i := l.Find(r)
                                if i >= 0 {
                                        found++
                                }
                        }
                        wg.Done()
                        log("Thread exited, found: %d", found)
                        return
                }()
        }

        // main loop
        t0 := time.Now()
        for {
                i, err := s.Generate()
                if err != nil {
                        break
                }
                lb.Assign(i)
        }
        // close all channels
        for j := 0; j < *k; j++ {
                close(cs[j])
        }
        // wait for all threads to finish
        wg.Wait()
        t1 := time.Now()
        d := t1.Sub(t0)
        fmt.Printf("%d\t%d\t%d\t%v\t%f\n", *k, *m, *n, d, float64(*n)/d.Seconds())
}
