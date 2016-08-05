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

type wrapper struct {
	bytes []byte
}

func newWrapper(bs []byte) *wrapper {
	w := new(wrapper)
	w.bytes = make([]byte, len(bs))
	copy(w.bytes, bs)
	return w
}

func (this *ListValue) ToRESP() respvalues.BasicType {
	elements := list.New()
	for elt := this.elements.Front(); elt != nil; elt = elt.Next() {
		val := elt.Value.(*wrapper).bytes
		elements.PushBack(respvalues.NewBulkStr(val))
	}
	return respvalues.NewArray(elements)
}

func (this *ListValue) PopFront() (value []byte, empty bool) {
	if this.elements.Len() == 0 {
		return nil, true
	}
	elt := this.elements.Front()
	this.elements.Remove(elt)

	value = elt.Value.(*wrapper).bytes
	empty = this.elements.Len() == 0

	return value, empty
}

func (this *ListValue) PopBack() (value []byte, empty bool) {
	if this.elements.Len() == 0 {
		return nil, true
	}
	elt := this.elements.Back()
	this.elements.Remove(elt)

	value = elt.Value.(*wrapper).bytes
	empty = this.elements.Len() == 0

	return value, empty
}

func (this *ListValue) PushBack(value []byte) int {
	wrapped := newWrapper(value)
	this.elements.PushBack(wrapped)
	return this.elements.Len()
}
func (this *ListValue) PushFront(value []byte) int {
	wrapped := newWrapper(value)
	this.elements.PushFront(wrapped)
	return this.elements.Len()
}

func (this *ListValue) Nth(idx int) ([]byte, bool) {
	if idx <= 0 || idx > this.elements.Len() {
		return nil, false
	}

	elt := this.elements.Front()
	for i := 1; i < idx; i++ {
		elt = elt.Next()
	}

	return elt.Value.(*wrapper).bytes, true
}
