package main

import (
	// "fmt"
	// "strings"

	// splay "github.com/gijsbers/go-splaytree"
	// splay "github.com/golang/collections/splay"
	splay "github.com/golang-collections/collections/splay"
)

// const SplaySparseness = 150

// type SplayItem struct {
//         i int
// }

// func (a SplayItem) Less(b splay.Item) bool {
//         return a.i < b.(SplayItem).i
// }

type SplayTree struct {
	tree *splay.SplayTree
}

func Less(a, b interface{}) bool {
	return a.(Item).Id() < b.(Item).Id()
}

func NewSplayTree() *SplayTree {
	log("SplayTree: %s", "creating")
	l := &SplayTree{tree: splay.New(Less)}
	return l
}

func (l *SplayTree) Add(i Item) {
	// log("SplayTree: adding item %d", i.Id())
	l.tree.Add(i)
}

func (l *SplayTree) String() string {
	return l.tree.String()
}

func (l *SplayTree) Find(val Item) int {
	// log("SplayTree: find id %d", val)
	i := l.tree.Get(val)
	if i != nil {
		return 1
	} else {
		return -1
	}
}
