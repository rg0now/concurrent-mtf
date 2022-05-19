package main

type IntegerItem struct {
        id Id
}

func (i IntegerItem) Id() Id {
        return i.id
}

func (i IntegerItem) Match(j Id) bool {
        return i.id == j
}
