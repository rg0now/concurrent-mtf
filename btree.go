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
        return l;
}

func (l *BTree) Add(i Item) {
        log("BTree: adding item %d", i.Id())
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
        log("BTree: find id %d", val)
        i := l.tree.Get(BTreeItem{i: val.Id()})
        if i != nil {
                return 1
        } else {
                return -1
        }
}
