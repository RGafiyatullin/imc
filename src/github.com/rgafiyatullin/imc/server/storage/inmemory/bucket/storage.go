package bucket

import (
	"errors"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
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

func (this *storage) handleCommand(cmd Cmd) (respvalues.BasicType, error) {
	switch cmd.CmdId() {
	case cmdGet:
		return this.handleCommandGet(cmd.(*CmdGet))
	case cmdSet:
		return this.handleCommandSet(cmd.(*CmdSet))
	case cmdExists:
		return this.handleCommandExists(cmd.(*CmdExists))
	case cmdDel:
		return this.handleCommandDel(cmd.(*CmdDel))
	case cmdLPushBack:
		return this.handleCommandLPushBack(cmd.(*CmdLPushBack))
	case cmdLPushFront:
		return this.handleCommandLPushFront(cmd.(*CmdLPushFront))
	case cmdLPopBack:
		return this.handleCommandLPopBack(cmd.(*CmdLPopBack))
	case cmdLPopFront:
		return this.handleCommandLPopFront(cmd.(*CmdLPopFront))
	case cmdLGetNth:
		return this.handleCommandLGetNth(cmd.(*CmdLGetNth))
	default:
		return nil, errors.New("unsupported command")
	}
}

func (this *storage) PurgeTimedOut() {
	for {
		key, exists := this.ttl.FetchTimedOut(this.tickIdx)
		if !exists {
			break
		}

		this.actorCtx.Log().Debug("purging timed out key '%s'", key)
		this.kv.Del(key)
	}
}

func (this *storage) handleCommandGet(cmd *CmdGet) (respvalues.BasicType, error) {
	kve, found := this.kv.Get(cmd.key)
	if !found {
		return respvalues.NewNil(), nil
	}

	validThru := kve.validThru()
	if validThru != 0 && validThru < this.tickIdx {
		this.kv.Del(cmd.key)
		this.ttl.SetTTL(cmd.key, 0)
		return respvalues.NewNil(), nil
	}

	return kve.value().ToRESP(), nil
}

func (this *storage) handleCommandSet(cmd *CmdSet) (respvalues.BasicType, error) {
	this.kv.Set(cmd.key, data.NewScalar(cmd.value), 0)
	this.ttl.SetTTL(cmd.key, 0)

	return respvalues.NewStr("OK"), nil
}

func (this *storage) handleCommandExists(cmd *CmdExists) (respvalues.BasicType, error) {
	return nil, errors.New("EXISTS: Not implemented")
}

func (this *storage) handleCommandDel(cmd *CmdDel) (respvalues.BasicType, error) {
	kve, existed := this.kv.Get(cmd.key)
	this.kv.Del(cmd.key)
	this.ttl.SetTTL(cmd.key, uint64(0))

	affectedRecords := int64(0)
	if existed && kve.validThru() >= this.tickIdx {
		affectedRecords = 1
	}

	return respvalues.NewInt(affectedRecords), nil
}

func (this *storage) handleCommandLPushBack(cmd *CmdLPushBack) (respvalues.BasicType, error) {
	return this.handleCommandLPushCommon(cmd.key, cmd.value, false)
}

func (this *storage) handleCommandLPushFront(cmd *CmdLPushFront) (respvalues.BasicType, error) {
	return this.handleCommandLPushCommon(cmd.key, cmd.value, true)
}

func (this *storage) handleCommandLPushCommon(key string, value []byte, front bool) (respvalues.BasicType, error) {
	kve, found := this.kv.Get(key)
	if found {
		switch kve.value().(type) {
		case (*data.ListValue):
			l := kve.value().(*data.ListValue)
			newlen := 0
			if front {
				newlen = l.PushFront(value)
			} else {
				newlen = l.PushBack(value)
			}

			return respvalues.NewInt(int64(newlen)), nil

		default:
			return respvalues.NewErr("LPSH*: incompatible existing value for this operation"), nil
		}
	} else {
		l := data.NewList()
		if front {
			l.PushFront(value)
		} else {
			l.PushBack(value)
		}
		this.kv.Set(key, l, 0)

		return respvalues.NewInt(int64(1)), nil
	}
}

func (this *storage) handleCommandLPopBack(cmd *CmdLPopBack) (respvalues.BasicType, error) {
	return nil, errors.New("LPOPB: not implemented [storage]")
}

func (this *storage) handleCommandLPopFront(cmd *CmdLPopFront) (respvalues.BasicType, error) {
	return nil, errors.New("LPOPF: not implemented [storage]")
}

func (this *storage) handleCommandLGetNth(cmd *CmdLGetNth) (respvalues.BasicType, error) {
	kve, found := this.kv.Get(cmd.key)
	if !found {
		return respvalues.NewNil(), nil
	}

	switch kve.value().(type) {
	case (*data.ListValue):
		l := kve.value().(*data.ListValue)
		value, found := l.Nth(cmd.idx)

		if !found {
			return respvalues.NewNil(), nil
		} else {
			return respvalues.NewBulkStr(value), nil
		}
	default:
		return respvalues.NewErr("LNTH: incompatible existing value for this operation"), nil
	}
}
