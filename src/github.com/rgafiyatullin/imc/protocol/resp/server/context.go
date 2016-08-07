package server

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/rgafiyatullin/imc/protocol/resp/constants"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"net"
	"net/textproto"
	"strconv"
)

// RESP-protocol context (over a net.Conn)
type Context interface {
	// Create new protocol context from socket
	Read() (*respvalues.RESPArray, error)

	// Read next protocol command (blocking)
	Write(data respvalues.RESPValue)
}
type context struct {
	sock net.Conn
	text *textproto.Conn
}

func New(sock net.Conn) Context {
	ctx := new(context)
	ctx.sock = sock
	ctx.text = textproto.NewConn(sock)
	return ctx
}

func (this *context) Read() (*respvalues.RESPArray, error) {
	line, err := this.text.ReadLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return this.Read()
	}
	switch line[0] {
	case constants.PrefixArray:
		cnt, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("EXPECTED_ARRAY_SIZE[GOT:'%v']", line[1:]))
		}
		elements, err := this.processArray(cnt)
		if err != nil {
			return nil, err
		}
		return elements, nil
	default:
		return nil, errors.New(fmt.Sprintf("EXPECTED_ARRAY_MARKER[GOT:'%v']", line[0]))
	}
}

// Write protocol command
func (this *context) Write(data respvalues.RESPValue) {
	data.Write(this.text)
}

func (this *context) processArray(count64 int64) (*respvalues.RESPArray, error) {
	count := int(count64)
	elements := list.New()
	for i := 0; i < count; i++ {
		element, err := this.processCommandElement()
		if err != nil {
			return nil, err
		}

		elements.PushBack(element)
	}
	return respvalues.NewArray(elements), nil
}

func (this *context) processCommandElement() (respvalues.RESPValue, error) {
	line, err := this.text.ReadLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return this.processCommandElement()
	}

	switch line[0] {
	case constants.PrefixInteger:
		value, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return nil, err
		}
		return respvalues.NewInt(value), nil

	case constants.PrefixError:
		return respvalues.NewErr(line[1:]), nil

	case constants.PrefixStr:
		return respvalues.NewStr(line[1:]), nil

	case constants.PrefixArray:
		count, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return nil, err
		}
		return this.processArray(count)

	case constants.PrefixBulkStr:
		bytesCount64, err := strconv.ParseInt(line[1:], 10, 64)
		bytesCount := int(bytesCount64)
		if err != nil {
			return nil, err
		}
		bytesPeeked, err := this.text.R.Peek(bytesCount)
		if err != nil {
			return nil, err
		}
		_, err = this.text.R.Discard(bytesCount)
		if err != nil {
			return nil, err
		}
		return respvalues.NewBulkStr(bytesPeeked), nil
	default:
		return nil, errors.New(fmt.Sprintf("UNEXPECTED_TYPE_PREFIX['%v']", line[0]))
	}
}
