package respvalues

import (
	"net/textproto"
)

// A common interface for the RESP-types (http://redis.io/topics/protocol#resp-protocol-description)
// Implemented by:
//
// * RESPArr
//
// * RESPBulkStr
//
// * RESPErr
//
// * RESPInt
//
// * RESPNil
//
// * RESPStr
type RESPValue interface {
	ToString() string
	Write(to *textproto.Conn)
}
