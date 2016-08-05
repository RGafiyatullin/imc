package bucket

import (
	"errors"
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
)

type storage struct {
	actorCtx actor.Ctx
	kv       KV
	ttl      TTL
}

func NewStorage() *storage {
	s := new(storage)
	s.kv = NewKV()
	s.ttl = NewTTL()
	return s
}

func (this *storage) handleCommand(cmd Cmd) (types.BasicType, error) {
	switch cmd.CmdId() {
	case cmdGet:
		return this.handleCommandGet(cmd.(*CmdGet))
	case cmdSet:
		return this.handleCommandSet(cmd.(*CmdSet))
	case cmdExists:
		return this.handleCommandExists(cmd.(*CmdExists))
	case cmdDel:
		return this.handleCommandDel(cmd.(*CmdDel))
	default:
		return nil, errors.New("unsupported command")
	}
}

func (this *storage) handleCommandGet(cmd *CmdGet) (types.BasicType, error) {
	return nil, errors.New("GET: Not implemented")
}

func (this *storage) handleCommandSet(cmd *CmdSet) (types.BasicType, error) {
	return nil, errors.New("SET: Not implemented")
}

func (this *storage) handleCommandExists(cmd *CmdExists) (types.BasicType, error) {
	return nil, errors.New("EXISTS: Not implemented")
}

func (this *storage) handleCommandDel(cmd *CmdDel) (types.BasicType, error) {
	return nil, errors.New("DEL: Not implemented")
}
