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

func (this *ListValue) ToRESP() respvalues.RESPValue {
	elements := list.New()
	for elt := this.elements.Front(); elt != nil; elt = elt.Next() {
		val := elt.Value.(*ScalarValue).value
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

	value = elt.Value.(*ScalarValue).value
	empty = this.elements.Len() == 0

	return value, empty
}

func (this *ListValue) PopBack() (value []byte, empty bool) {
	if this.elements.Len() == 0 {
		return nil, true
	}
	elt := this.elements.Back()
	this.elements.Remove(elt)

	value = elt.Value.(*ScalarValue).value
	empty = this.elements.Len() == 0

	return value, empty
}

func (this *ListValue) PushBack(value []byte) int {
	wrapped := NewScalar(value)
	this.elements.PushBack(wrapped)
	return this.elements.Len()
}
func (this *ListValue) PushFront(value []byte) int {
	wrapped := NewScalar(value)
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

	return elt.Value.(*ScalarValue).value, true
}
