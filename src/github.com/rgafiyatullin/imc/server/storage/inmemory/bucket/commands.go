package bucket

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
)

const cmdSet = 1
const cmdGet = 2
const cmdExists = 3
const cmdDel = 4

type Cmd interface {
	CmdId() int
}

func NewCmdGet(key string) Cmd {
	cmd := new(CmdGet)
	cmd.key = key
	return cmd
}

func NewCmdSet(key string, value respvalues.BasicType, expiry uint64) Cmd {
	cmd := new(CmdSet)
	cmd.key = key
	cmd.value = value
	cmd.expiry = expiry
	return cmd
}

func NewCmdDel(key string) Cmd {
	cmd := new(CmdDel)
	cmd.key = key
	return cmd
}

type CmdSet struct {
	key    string
	value  respvalues.BasicType
	expiry uint64
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
