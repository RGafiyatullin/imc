package types

import (
	"fmt"
	"net/textproto"
)

func NewBulkStr(v []byte) *BasicBulkStr {
	s := new(BasicBulkStr)
	s.bytes = v
	return s
}

type BasicBulkStr struct {
	bytes []byte
}

func (this *BasicBulkStr) Bytes() []byte {
	return this.bytes
}
func (this *BasicBulkStr) String() string {
	return fmt.Sprintf("%s", this.bytes)
}

func (this *BasicBulkStr) ToString() string {
	return fmt.Sprintf("B(\"%s\")", this.bytes)
}
func (this *BasicBulkStr) Write(to *textproto.Conn) {
	to.Cmd("$%d", len(this.bytes))
	to.W.Write(this.bytes)
	to.Cmd("")
}
