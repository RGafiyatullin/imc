package respvalues

import (
	"fmt"
	"net/textproto"
)

func NewErr(v string) *RESPErr {
	e := new(RESPErr)
	e.e = v
	return e
}

// Represents RESP-Error value (http://redis.io/topics/protocol#resp-errors)
type RESPErr struct {
	e string
}

func (this *RESPErr) ToString() string {
	return fmt.Sprintf("E(\"%s\")", this.e)
}
func (this *RESPErr) Write(to *textproto.Conn) {
	to.Cmd("-%s", this.e)
}
