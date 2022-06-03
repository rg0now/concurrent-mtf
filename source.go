package main

import (
        "fmt"
        "math/rand"

        "gonum.org/v1/gonum/stat/distuv"
)
type UniformSource struct {
        m, n, i int
}

func (s *UniformSource) Generate() (Item, error) {
        if s.i < s.n {
                s.i++
                i := rand.Intn(s.m)
                log("Source: generating req %d: %d", s.i, i)
                return IntegerItem{i}, nil
        } else {
                return IntegerItem{}, fmt.Errorf("Done")
        }
}

type PoissonSource struct {
        m, n, i int
        rand *distuv.Poisson
        store []Item
}


func NewPoissonSource(m, n int, la float64, store *([]Item)) *PoissonSource {
        s := &PoissonSource{m: m, n: n, i: 0, rand: &distuv.Poisson{Lambda: la},
                store: make([]Item, m)}
        copy(s.store, *store)
        rand.Shuffle(m, func(i, j int) {
		s.store[i], s.store[j] = s.store[j], s.store[i]
	})
        log("%s", s.store)
        return s 
}

func (s *PoissonSource) Generate() (Item, error) {
        if s.i < s.n {
                s.i++
                i := int(s.rand.Rand())
                if i >= s.m {
                        i = s.m-1
                }
                log("Source: generating req %d: %d", s.i, s.store[i].Id())
                return s.store[i], nil
        } else {
                return IntegerItem{}, fmt.Errorf("Done")
        }
}
