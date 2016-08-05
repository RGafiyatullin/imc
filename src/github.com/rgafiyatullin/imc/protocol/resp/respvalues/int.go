package respvalues

import (
	"fmt"
	"net/textproto"
)

func NewInt(v int64) *BasicInt {
	i := new(BasicInt)
	i.i = v
	return i
}

type BasicInt struct {
	i int64
}

func (this *BasicInt) ToString() string {
	return fmt.Sprintf("\"I(%d)\"", this.i)
}
func (this *BasicInt) Write(to *textproto.Conn) {
	to.Cmd(":%d", this.i)
}
func (this *BasicInt) Plus(other *BasicInt) *BasicInt {
	return NewInt(this.i + other.i)
}
