package main

import (
	// "fmt"
	// "strings"

	"github.com/google/btree"
)

type BTreeItem struct {
	i int
}

func (a BTreeItem) Less(b btree.Item) bool {
	return a.i < b.(BTreeItem).i
}

type BTree struct {
	tree *btree.BTree
}

func NewBTree() *BTree {
	log("BTree: %s", "creating")
	l := &BTree{tree: btree.New(2)}
	return l
}

func (l *BTree) Add(i Item) {
	// log("BTree: adding item %d", i.Id())
	l.tree.ReplaceOrInsert(BTreeItem{i: i.Id()})
}

func (l *BTree) String() string {
	// var ns []string
	// size := 0
	// l.tree.Ascend(func(a Item) bool {
	// 	ns = append(ns, fmt.Sprintf("%v ", a))
	//         size += 1
	//         return true
	// })
	// return fmt.Sprint("BTree: size=%d: ") + strings.Join(ns, ", ")
	return "BTree"
}

func (l *BTree) Find(val Item) int {
	// log("BTree: find id %d", val)
	i := l.tree.Get(BTreeItem{i: val.Id()})
	if i != nil {
		return 1
	} else {
		return -1
	}
}

// WeightedBTree is a balanced tree with a somewhat heavier comparison operator
type WeightedBTreeItem struct {
	i int
}

func (a WeightedBTreeItem) Less(b btree.Item) bool {
	sum := 0
	for i := 0; i <= WeightedTreeBusyWait; i++ {
		sum += i
	}
	if sum/2 == 0 {
		return a.i < b.(WeightedBTreeItem).i
	} else {
		return b.(WeightedBTreeItem).i > a.i
	}
}

type WeightedBTree struct {
	tree *btree.BTree
}

func NewWeightedBTree() *WeightedBTree {
	log("Weighted BTree: creating, busyweight in less: %d", WeightedTreeBusyWait)
	l := &WeightedBTree{tree: btree.New(2)}
	return l
}

func (l *WeightedBTree) Add(i Item) {
	// log("BTree: adding item %d", i.Id())
	l.tree.ReplaceOrInsert(WeightedBTreeItem{i: i.Id()})
}

func (l *WeightedBTree) String() string {
	return "WeightedBTree"
}

func (l *WeightedBTree) Find(val Item) int {
	// log("BTree: find id %d", val)
	i := l.tree.Get(WeightedBTreeItem{i: val.Id()})
	if i != nil {
		return 1
	} else {
		return -1
	}
}
