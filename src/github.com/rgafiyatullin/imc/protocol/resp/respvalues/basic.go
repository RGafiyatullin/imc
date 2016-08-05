package respvalues

import (
	"net/textproto"
)

type BasicType interface {
	ToString() string
	Write(to *textproto.Conn)
}
