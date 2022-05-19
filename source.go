package main

import (
        "fmt"
        "math/rand"

        "gonum.org/v1/gonum/stat/distuv"
)
type UniformSource struct {
        m, n, i int
}

func (s *UniformSource) Generate() (Id, error) {
        if s.i < s.n {
                s.i++
                i := rand.Intn(s.m)
                log("Source: generating req %d: %d", s.i, i)
                return i, nil
        } else {
                return 0, fmt.Errorf("Done")
        }
}

type PoissonSource struct {
        m, n, i int
        rand *distuv.Poisson
}


func NewPoissonSource(m, n int, la float64) *PoissonSource {
        return &PoissonSource{m: m, n: n, i: 0, rand: &distuv.Poisson{Lambda: la}}
}

func (s *PoissonSource) Generate() (Id, error) {
        if s.i < s.n {
                s.i++
                i := int(s.rand.Rand())
                if i > s.m {
                        i = s.m
                }
                log("Source: generating req %d: %d", s.i, i)
                return i, nil
        } else {
                return 0, fmt.Errorf("Done")
        }
}
