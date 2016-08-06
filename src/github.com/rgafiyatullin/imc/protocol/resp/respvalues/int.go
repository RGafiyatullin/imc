package respvalues

import (
	"fmt"
	"net/textproto"
)

func NewInt(v int64) *RESPInt {
	i := new(RESPInt)
	i.i = v
	return i
}

// Represents RESP-Int value (http://redis.io/topics/protocol#resp-integers)
type RESPInt struct {
	i int64
}

func (this *RESPInt) ToString() string {
	return fmt.Sprintf("\"I(%d)\"", this.i)
}
func (this *RESPInt) Write(to *textproto.Conn) {
	to.Cmd(":%d", this.i)
}
func (this *RESPInt) Plus(other *RESPInt) *RESPInt {
	return NewInt(this.i + other.i)
}
