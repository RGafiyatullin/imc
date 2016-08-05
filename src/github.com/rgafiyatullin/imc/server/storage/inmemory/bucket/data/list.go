package data

import (
	"container/list"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
)

type ListValue struct {
	elements *list.List
}

func NewList() *ListValue {
	l := new(ListValue)
	l.elements = list.New()
	return l
}

func (this *ListValue) ToRESP() respvalues.BasicType {
	elements := list.New()
	for elt := this.elements.Front(); elt != nil; elt = elt.Next() {
		val := elt.Value.([]byte)
		elements.PushBack(respvalues.NewBulkStr(val))
	}
	return respvalues.NewArray(elements)
}

func (this *ListValue) Append(value []byte) int {
	this.elements.PushBack(value)
	return this.elements.Len()
}
func (this *ListValue) Prepend(value []byte) int {
	this.elements.PushFront(value)
	return this.elements.Len()
}
