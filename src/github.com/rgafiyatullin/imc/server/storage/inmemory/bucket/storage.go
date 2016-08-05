package bucket

import (
	"errors"
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
)

type storage struct {
	actorCtx actor.Ctx
	tickIdx  uint64
	kv       KV
	ttl      TTL
}

func NewStorage() *storage {
	s := new(storage)
	s.tickIdx = 0
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
	kve, found := this.kv.Get(cmd.key)
	if !found {
		return types.NewNil(), nil
	}
	if kve.validThru() < this.tickIdx && kve.validThru() != 0 {
		this.kv.Del(cmd.key)
		this.ttl.SetTTL(cmd.key, 0)
		return types.NewNil(), nil
	}

	return kve.value(), nil
}

func (this *storage) handleCommandSet(cmd *CmdSet) (types.BasicType, error) {
	validThru := this.tickIdx + cmd.expiry
	if cmd.expiry == 0 {
		validThru = 0
	}

	this.kv.Set(cmd.key, cmd.value, validThru)
	this.ttl.SetTTL(cmd.key, validThru)

	return types.NewStr("OK"), nil
}

func (this *storage) handleCommandExists(cmd *CmdExists) (types.BasicType, error) {
	return nil, errors.New("EXISTS: Not implemented")
}

func (this *storage) handleCommandDel(cmd *CmdDel) (types.BasicType, error) {
	kve, existed := this.kv.Get(cmd.key)
	this.kv.Del(cmd.key)
	this.ttl.SetTTL(cmd.key, uint64(0))

	affectedRecords := int64(0)
	if existed && kve.validThru() >= this.tickIdx {
		affectedRecords = 1
	}

	return types.NewInt(affectedRecords), nil
}
