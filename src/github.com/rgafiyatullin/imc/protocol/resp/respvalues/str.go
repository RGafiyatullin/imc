package respvalues

import (
	"fmt"
	"net/textproto"
)

func NewStr(v string) *RESPStr {
	s := new(RESPStr)
	s.s = v
	return s
}

// Represents RESP-SimpleString value (http://redis.io/topics/protocol#resp-simple-strings)
type RESPStr struct {
	s string
}

func (this *RESPStr) ToString() string {
	return fmt.Sprintf("S(\"%s\")", this.s)
}
func (this *RESPStr) Write(to *textproto.Conn) {
	to.Cmd("+%s", this.s)
}
