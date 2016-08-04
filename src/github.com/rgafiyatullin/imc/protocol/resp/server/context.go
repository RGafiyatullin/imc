package server

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/rgafiyatullin/imc/protocol/resp/constants"
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"net"
	"net/textproto"
	"strconv"
)

type Context interface {
	NextCommand() (*types.BasicArr, error)
	Write(data types.BasicType)
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

func (this *context) NextCommand() (*types.BasicArr, error) {
	line, err := this.text.ReadLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return this.NextCommand()
	}
	switch line[0] {
	case constants.PrefixArray:
		cnt, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("expected array size, got '%v'", line[1:]))
		}
		elements, err := this.processArray(cnt)
		if err != nil {
			return nil, err
		}
		return elements, nil
	default:
		return nil, errors.New(fmt.Sprintf("expected array (*), got '%v'", line[0]))
	}
}

func (this *context) Write(data types.BasicType) {
	data.Write(this.text)
}

func (this *context) processArray(count64 int64) (*types.BasicArr, error) {
	count := int(count64)
	elements := list.New()
	for i := 0; i < count; i++ {
		element, err := this.processCommandElement()
		if err != nil {
			return nil, err
		}

		elements.PushBack(element)
	}
	return types.NewArray(elements), nil
}

func (this *context) processCommandElement() (types.BasicType, error) {
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
		return types.NewInt(value), nil

	case constants.PrefixError:
		return types.NewErr(line[1:]), nil

	case constants.PrefixStr:
		return types.NewStr(line[1:]), nil

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
		return types.NewBulkStr(bytesPeeked), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unexpected type-prefix: '%v'", line[0]))
	}
}
