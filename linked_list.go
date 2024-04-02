package main

import (
	"fmt"
	"strings"
)

type LinkedListNode struct {
	value Item
	next  *LinkedListNode
}

type LinkedList struct {
	head *LinkedListNode
	len  int
}

func NewLinkedList() *LinkedList {
	log("LinkedList: creating")
	return &LinkedList{head: nil, len: 0}
}

func (l *LinkedList) Add(i Item) {
	// log("LinkedList: adding item %d", i.Id())
	n := LinkedListNode{}
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
func (l *LinkedList) String() string {
	ret := fmt.Sprintf("LinkedList: size %d: ", l.len)
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
func (l *LinkedList) Find(val Item) int {
	// log("LinkedList: find id %d", val.Id())
	ptr := l.head
	for j := 0; j < l.len; j++ {
		if ptr.value.Match(val.Id()) {
			if verbose {
				// log("LinkedList: found id %d at position %d", ptr.value, j)
				// log(l.String())
			}
			return j
		}
		ptr = ptr.next
	}
	return -1
}
