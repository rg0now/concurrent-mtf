package main

import (
        "fmt"        
        "strings"

        splay "github.com/gijsbers/go-splaytree"
)

// const SplaySparseness = 150

type SplayTree struct {
        tree *splay.SplayTree 
}

func (i IntegerItem) Less(than splay.Item) bool {
	return i.id < than.(IntegerItem).id
}

func NewSplayTree() *SplayTree {
        log("SplayTree: %s", "creating")
        l := &SplayTree{tree: splay.NewSplayTree()}
        return l;
}

func (l *SplayTree) Add(i *Item) {
        log("SplayTree: adding item %d", (*i).Id())
        l.tree.Insert((*i).(IntegerItem))
}

func (l *SplayTree) String() string {
        ret := fmt.Sprintf("SplayTree: size %d, root: %d: ", l.tree.Count(), l.tree.Root())
        var ns []string
	l.tree.Traverse(func(i splay.Item) { ns = append(ns, fmt.Sprintf("%v ", i)) })
        return ret + strings.Join(ns, ", ")
}

func (l *SplayTree) Find(val Id) int {
        log("SplayTree: find id %d", val)
        i := l.tree.Lookup(IntegerItem{val})
        if i != nil {
                return 1
        } else {
                return -1
        }
}
