package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat/distuv"
)

func shuffle(to, from *[]Item) {
	var rnd = rand.New(rand.NewSource(seed))
	*to = make([]Item, len(*from))
	copy(*to, *from)
	rnd.Shuffle(len(*to), func(i, j int) {
		(*to)[i], (*to)[j] = (*to)[j], (*to)[i]
	})
}

// UNIFORM
type UniformSource struct {
	m, n, i int
	rnd     *rand.Rand
	store   []Item
	offset  int
}

func NewUniformSource(m, n int, store *[]Item) *UniformSource {
	s := UniformSource{n: n, m: m, i: 0, rnd: rand.New(rand.NewSource(seed)), offset: 1.0}
	s.offset = 0
	if isBloom {
		// bloom and sabloom can process requests from a larger universe
		s.offset = int((1.0 - BloomFillRatio) * float64(s.m))
	}
	shuffle(&s.store, store)
	return &s
}

func (s *UniformSource) Generate() (Item, error) {
	if s.i < s.n {
		s.i++
		i := s.rnd.Intn(s.m)
		val := s.store[i].Id() + s.offset
		log("Source: generating req %d: %d, store: %d, offset: %d", s.i, val, s.store[i].Id(), s.offset)
		return IntegerItem{val}, nil
	} else {
		return IntegerItem{}, fmt.Errorf("Done")
	}
}

// POISSON
type PoissonSource struct {
	m, n, i int
	rand    *distuv.Poisson
	store   []Item
}

func NewPoissonSource(m, n int, la float64, store *([]Item)) *PoissonSource {
	// TODO: use seed
	s := &PoissonSource{m: m, n: n, i: 0, rand: &distuv.Poisson{Lambda: la}}
	shuffle(&s.store, store)
	// log("%s", s.store)
	return s
}

func (s *PoissonSource) Generate() (Item, error) {
	if s.i < s.n {
		s.i++
		i := int(s.rand.Rand())
		if i >= s.m {
			i = s.m - 1
		}
		// log("Source: generating req %d: %d", s.i, s.store[i].Id())
		return s.store[i], nil
	} else {
		return IntegerItem{}, fmt.Errorf("Done")
	}
}

// ZIPF
type ZipfSource struct {
	*rand.Zipf
	m, n, i int
	store   []Item
}

func NewZipfSource(name string, m, n int, store *([]Item)) *ZipfSource {
	rnd := rand.New(rand.NewSource(seed))
	s := 1.01
	args := strings.Split(name, ":")
	if len(args) == 2 {
		if f, err := strconv.ParseFloat(args[1], 64); err == nil {
			s = f
		}
	}
	z := rand.NewZipf(rnd, s, 1.0, uint64(m)-1)
	if z == nil {
		panic("cannot create zipf distribution")
	}
	zs := &ZipfSource{Zipf: z, m: m, n: n, i: 0}
	shuffle(&zs.store, store)
	// log("%s", s.store)
	return zs
}

func (s *ZipfSource) Generate() (Item, error) {
	if s.i < s.n {
		s.i++
		i := int(s.Uint64())
		log("-------%d", s.store[i])
		return s.store[i], nil
	} else {
		return IntegerItem{}, fmt.Errorf("Done")
	}
}
