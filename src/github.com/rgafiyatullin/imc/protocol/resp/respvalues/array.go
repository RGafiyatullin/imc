package respvalues

import (
	"container/list"
	"net/textproto"
)

func NewArray(el *list.List) *RESPArray {
	a := new(RESPArray)
	ea := make([]RESPValue, el.Len())
	i := 0
	for e := el.Front(); e != nil; e = e.Next() {
		ea[i] = e.Value.(RESPValue)
		i++
	}
	a.elements = ea
	return a
}

// Represents RESP-array value (http://redis.io/topics/protocol#resp-arrays)
type RESPArray struct {
	elements []RESPValue
}

func (this *RESPArray) Elements() []RESPValue {
	return this.elements
}

func (this *RESPArray) ToString() string {
	acc := "A("
	for i := 0; i < len(this.elements); i++ {
		acc += this.elements[i].ToString() + ", "
	}
	acc += ")"
	return acc
}

func (this *RESPArray) Write(to *textproto.Conn) {
	to.Cmd("*%d", len(this.elements))
	for i := 0; i < len(this.elements); i++ {
		this.elements[i].Write(to)
	}
}
