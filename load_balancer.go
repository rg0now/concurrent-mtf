package main

type ModuloLB struct {
	k    int
	s    Source
	pool []Channel
}

func NewModuloLB(k int, s Source, cs []Channel) *ModuloLB {
	log("ModuloLB: Creating ModuloLB for %d threads", k)
	lb := &ModuloLB{k: k, s: s, pool: cs}
	return lb
}

func (lb *ModuloLB) Generate() error {
	i, err := lb.s.Generate()
	if err != nil {
		return err
	}
	lb.Assign(i)
	return nil
}

func (lb *ModuloLB) Assign(r Item) {
	j := r.Id() % lb.k
	// log("ModuloLB: Assigning item %d to thread %d", r.Id(), j)
	lb.pool[j] <- r
}

type SplitLB struct {
	k, m int
	s    Source
	pool []Channel
}

func NewSplitLB(k, m int, s Source, cs []Channel) *SplitLB {
	log("SplitLB: Creating SplitLB over range [0,%d] for %d threads", k, m)
	lb := &SplitLB{m: m, k: k, s: s, pool: cs}
	return lb
}

func (lb *SplitLB) Generate() error {
	i, err := lb.s.Generate()
	if err != nil {
		return err
	}
	lb.Assign(i)
	return nil
}

func (lb *SplitLB) Assign(r Item) {
	j := int(float64(r.Id()) / (float64(lb.m) / float64(lb.k)))
	// log("SplitLB: Assigning item %d to thread %d", r, j)
	lb.pool[j] <- r
}

type RoundRobinLB struct {
	k, i int
	s    Source
	pool []Channel
}

func NewRoundRobinLB(k int, s Source, cs []Channel) *RoundRobinLB {
	log("RoundRobinLB: Creating RoundRobinLB for %d threads", k)
	lb := &RoundRobinLB{i: 0, k: k, s: s, pool: cs}
	return lb
}

func (lb *RoundRobinLB) Generate() error {
	i, err := lb.s.Generate()
	if err != nil {
		return err
	}
	lb.Assign(i)
	return nil
}

func (lb *RoundRobinLB) Assign(r Item) {
	j := lb.i % lb.k
	// log("RoundRobinLB: Assigning item %d to thread %d", r.Id(), j)
	lb.pool[j] <- r
	lb.i += 1
}
