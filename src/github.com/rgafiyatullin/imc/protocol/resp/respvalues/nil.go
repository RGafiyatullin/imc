package respvalues

import (
	"net/textproto"
)

func NewNil() *RESPNil {
	n := new(RESPNil)
	return n
}

// See http://redis.io/topics/protocol#resp-bulk-strings (Null Bulk String)
type RESPNil struct{}

func (this *RESPNil) ToString() string {
	return "NIL"
}
func (this *RESPNil) Write(to *textproto.Conn) {
	to.Cmd("$-1")
}
