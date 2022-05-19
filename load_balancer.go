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

func (lb *ModuloLB) Assign(r Id) {
        j := r % lb.k
        log("ModuloLB: Assigning item %d to thread %d", r, j)
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

func (lb *SplitLB) Assign(r Id) {
        j := int(float64(r) /  (float64(lb.m) / float64(lb.k)))
        log("SplitLB: Assigning item %d to thread %d", r, j)
        lb.pool[j] <- r
}
