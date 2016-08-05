package bucket

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"time"
)

const cmdSet = 1
const cmdGet = 2
const cmdExists = 3
const cmdDel = 4

type Cmd interface {
	CmdId() int
}

type CmdSet struct {
	key   string
	value types.BasicType
	ttl   time.Duration
}

func (this *CmdSet) CmdId() int {
	return cmdSet
}

type CmdGet struct {
	key string
}

func (this *CmdGet) CmdId() int {
	return cmdGet
}

type CmdExists struct {
	key string
}

func (this *CmdExists) CmdId() int {
	return cmdExists
}

type CmdDel struct {
	key string
}

func (this *CmdDel) CmdId() int {
	return cmdDel
}
