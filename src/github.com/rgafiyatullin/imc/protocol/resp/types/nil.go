package types

import (
	"net/textproto"
)

func NewNil() *BasicNil {
	n := new(BasicNil)
	return n
}

type BasicNil struct{}

func (this *BasicNil) ToString() string {
	return "NIL"
}
func (this *BasicNil) Write(to *textproto.Conn) {
	to.Cmd("$-1")
}
