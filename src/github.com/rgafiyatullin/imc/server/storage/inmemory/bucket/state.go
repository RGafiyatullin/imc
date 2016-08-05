package bucket

import (
	"errors"
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
)

type state struct {
	actorCtx actor.Ctx
	kv       KV
}

func NewState(actorCtx actor.Ctx) state {
	s := new(state)
	s.actorCtx = actorCtx
	s.kv = NewKV()
	return s
}

func (this *state) handleCommand(cmd Cmd) (*types.BasicType, error) {
	switch cmd.CmdId() {
	case cmdGet:
		return this.handleCommandGet(cmd.(CmdGet))
	case cmdSet:
		return this.handleCommandSet(cmd.(CmdSet))
	case cmdExists:
		return this.handleCommandExists(cmd.(CmdExists))
	case cmdDel:
		return this.handleCommandDel(cmd.(CmdDel))
	default:
		return nil, errors.New("unsupported command")
	}
}

func (this *state) handleCommandGet(cmd *CmdGet) (*types.BasicType, error) {
	return nil, errors.New("GET: Not implemented")
}

func (this *state) handleCommandSet(cmd *CmdSet) (*types.BasicType, error) {
	return nil, errors.New("SET: Not implemented")
}

func (this *state) handleCommandExists(cmd *CmdExists) (*types.BasicType, error) {
	return nil, errors.New("EXISTS: Not implemented")
}

func (this *state) handleCommandDel(cmd *CmdDel) (*types.BasicType, error) {
	return nil, errors.New("DEL: Not implemented")
}
