package types

import (
	"fmt"
	"net/textproto"
)

func NewErr(v string) *BasicErr {
	e := new(BasicErr)
	e.e = v
	return e
}

type BasicErr struct {
	e string
}

func (this *BasicErr) ToString() string {
	return fmt.Sprintf("E(\"%s\")", this.e)
}
func (this *BasicErr) Write(to *textproto.Conn) {
	to.Cmd("-%s", this.e)
}
