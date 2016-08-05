package types

import (
	"fmt"
	"net/textproto"
)

func NewStr(v string) *BasicStr {
	s := new(BasicStr)
	s.s = v
	return s
}

type BasicStr struct {
	s string
}

func (this *BasicStr) ToString() string {
	return fmt.Sprintf("S(\"%s\")", this.s)
}
func (this *BasicStr) Write(to *textproto.Conn) {
	to.Cmd("+%s", this.s)
}
