package main

type ModuloLB struct {
        k int
        pool []Channel
}

func NewModuloLB(k int, cs []Channel) *ModuloLB {
        log("ModuloLB: Creating ModuloLB for %d threads", k)
        lb := &ModuloLB{k: k, pool: cs}
        return lb
}

func (lb *ModuloLB) Assign(r Item) {
        j := r.Id() % lb.k
        log("ModuloLB: Assigning item %d to thread %d", r.Id(), j)
        lb.pool[j] <- r
}

type SplitLB struct {
        k, m int
        pool []Channel
}

func NewSplitLB(k, m int, cs []Channel) *SplitLB {
        log("SplitLB: Creating SplitLB over range [0,%d] for %d threads", k, m)
        lb := &SplitLB{m: m, k: k, pool: cs}
        return lb
}

func (lb *SplitLB) Assign(r Item) {
        j := int(float64(r.Id()) /  (float64(lb.m) / float64(lb.k)))
        log("SplitLB: Assigning item %d to thread %d", r, j)
        lb.pool[j] <- r
}

type RoundRobinLB struct {
        k, i int
        pool []Channel
}

func NewRoundRobinLB(k int, cs []Channel) *RoundRobinLB {
        log("RoundRobinLB: Creating RoundRobinLB for %d threads", k)
        lb := &RoundRobinLB{i: 0, k: k, pool: cs}
        return lb
}

func (lb *RoundRobinLB) Assign(r Item) {
        j := lb.i % lb.k
        log("RoundRobinLB: Assigning item %d to thread %d", r.Id(), j)
        lb.pool[j] <- r
        lb.i += 1
}
