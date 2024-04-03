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

// WeightedSplay tree is a splay-tree implementation with a somewhat more heavy-weight less
// function for comparison. This is needed orterwise LB becomes the bottleneck and this is not
// something we're interested in for now
type WeightedSplayTree struct {
	SplayTree
}

func WeightedLess(a, b interface{}) bool {
	sum := 0
	for i := 0; i <= WeightedTreeBusyWait; i++ {
		sum += i
	}
	if sum/2 == 0 {
		return a.(Item).Id() < b.(Item).Id()
	} else {
		return b.(Item).Id() > a.(Item).Id()
	}
}

func NewWeightedSplayTree() *WeightedSplayTree {
	log("Weighted SplayTree: creating, busyweight in less: %d", WeightedTreeBusyWait)
	l := &WeightedSplayTree{SplayTree: SplayTree{tree: splay.New(WeightedLess)}}
	return l
}
