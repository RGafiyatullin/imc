package respvalues

import (
	"fmt"
	"net/textproto"
)

func NewBulkStr(v []byte) *RESPBulkStr {
	s := new(RESPBulkStr)
	s.bytes = v
	return s
}

// Represents RESP-BulkString value (http://redis.io/topics/protocol#resp-bulk-strings)
type RESPBulkStr struct {
	bytes []byte
}

func (this *RESPBulkStr) Bytes() []byte {
	return this.bytes
}
func (this *RESPBulkStr) String() string {
	return fmt.Sprintf("%s", this.bytes)
}

func (this *RESPBulkStr) ToString() string {
	return fmt.Sprintf("B(\"%s\")", this.bytes)
}
func (this *RESPBulkStr) Write(to *textproto.Conn) {
	to.Cmd("$%d", len(this.bytes))
	to.W.Write(this.bytes)
	to.Cmd("")
}
