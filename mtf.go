package main

import (
        "fmt"
)

type MtfItem struct {
        id Id
}

func (i MtfItem) Id() Id {
        return i.id
}

func (i MtfItem) Match(j Id) bool {
        return i.id == j
}

type MtfNode struct {
        value *Item
        next  *MtfNode
}

type Mtf struct {
        head *MtfNode
        len int
}

func (l *Mtf) Add(i *Item) {
        log("Mtf: adding item %d", (*i).Id())
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
func (l *Mtf) Dump() {
	if l.len == 0 {
		fmt.Println("-")
	}
	ptr := l.head
        for j := 0; j < l.len; j++ {
		fmt.Println("Node: ", (*ptr.value).Id())
		ptr = ptr.next
	}
}

// Search returns node position with given value from linked list
func (l *Mtf) Find(val Id) int {
        log("Mtf: find id %d", val)
	ptr := l.head
        var prev *MtfNode = nil
        for j := 0; j < l.len; j++ {
		if (*ptr.value).Match(val) {
                        // MTF
                        // first item aleady?
                        if prev != nil {
                                // remove from list
                                prev.next = ptr.next
        
                                // move to head
                                ptr.next = l.head
                                l.head = ptr
                        }

                        if verbose {l.Dump()}
                        return j
		}
                prev = ptr
		ptr = ptr.next
	}
	return -1
}
