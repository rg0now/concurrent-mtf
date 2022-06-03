package main

import (
        "fmt"
        "strings"
)

type MtfNode struct {
        value Item
        next  *MtfNode
}

type Mtf struct {
        head *MtfNode
        len int
}

func NewMtf() *Mtf {
        return &Mtf{head: nil, len: 0}
}

func (l *Mtf) Add(i Item) {
        log("Mtf: adding item %d", i.Id())
        n := MtfNode{}
        n.value = i
        if l.len == 0 {
                l.head = &n
                l.len++
                return
        }
        ptr := l.head
        for j := 0; j < l.len; j++ {
                if ptr.next == nil {
                        ptr.next = &n
                        l.len++
                        return
                }
                ptr = ptr.next
        }
}

// Print displays all the nodes from linked list
func (l *Mtf) String() string {
        ret := fmt.Sprintf("Mtf: size %d: ", l.len)
        var ns []string
	if l.len == 0 {
		ns = append(ns, "-")
	} else {
                ptr := l.head
                for j := 0; j < l.len; j++ {
                        ns = append(ns, fmt.Sprintf("Node: %d", ptr.value.Id()))
                        ptr = ptr.next
                }
	}
        return ret + strings.Join(ns, ", ")       
}

// Search returns node position with given value from linked list
func (l *Mtf) Find(val Item) int {
        log("Mtf: find id %d", val)
	ptr := l.head
        var prev *MtfNode = nil
        for j := 0; j < l.len; j++ {
		if ptr.value.Match(val.Id()) {
                        // MTF
                        // first item aleady?
                        if prev != nil {
                                // remove from list
                                prev.next = ptr.next
        
                                // move to head
                                ptr.next = l.head
                                l.head = ptr
                        }

                        if verbose {
                                log("Mtf: found id %d at position %d", ptr.value, j)
                                log(l.String())
                        }
                        return j
		}
                prev = ptr
		ptr = ptr.next
	}
	return -1
}

