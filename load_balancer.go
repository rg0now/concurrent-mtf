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
