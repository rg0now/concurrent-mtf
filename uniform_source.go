package main

import (
        "fmt"
        "math/rand"
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

