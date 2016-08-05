package types

import (
	"container/list"
	"net/textproto"
)

func NewArray(el *list.List) *BasicArr {
	a := new(BasicArr)
	ea := make([]BasicType, el.Len())
	i := 0
	for e := el.Front(); e != nil; e = e.Next() {
		ea[i] = e.Value.(BasicType)
		i++
	}
	a.elements = ea
	return a
}

type BasicArr struct {
	elements []BasicType
}

func (this *BasicArr) ToString() string {
	acc := "A("
	for i := 0; i < len(this.elements); i++ {
		acc += this.elements[i].ToString() + ", "
	}
	acc += ")"
	return acc
}

func (this *BasicArr) Write(to *textproto.Conn) {
	to.Cmd("*%d", len(this.elements))
	for i := 0; i < len(this.elements); i++ {
		this.elements[i].Write(to)
	}
}
