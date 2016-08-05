package types

import (
	"fmt"
	"net/textproto"
)

func NewBulkStr(v []byte) *BasicBulkStr {
	s := new(BasicBulkStr)
	s.s = v
	return s
}

type BasicBulkStr struct {
	s []byte
}

func (this *BasicBulkStr) ToString() string {
	return fmt.Sprintf("B(\"%s\")", this.s)
}
func (this *BasicBulkStr) Write(to *textproto.Conn) {
	to.Cmd("$%d", len(this.s))
	to.W.Write(this.s)
	to.Cmd("")
}
