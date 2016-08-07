package bucket

import (
	"container/list"
	"errors"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/metronome"
	"github.com/steveyen/gtreap"
	"regexp"
	"time"
)

type storage struct {
	actorCtx actor.Ctx
	tickIdx  int64
	kv       KV
	ttl      TTL
}

func NewStorage(actorCtx actor.Ctx) *storage {
	s := new(storage)
	s.actorCtx = actorCtx
	s.tickIdx = 0
	s.kv = NewKV()
	s.ttl = NewTTL()
	return s
}

func (this *storage) handleCommand(cmd Cmd) (respvalues.RESPValue, error) {
	switch cmd.CmdId() {
	case cmdGet:
		return this.handleCommandGet(cmd.(*CmdGet))
	case cmdSet:
		return this.handleCommandSet(cmd.(*CmdSet))
	case cmdDel:
		return this.handleCommandDel(cmd.(*CmdDel))
	case cmdKeys:
		return this.handleCommandKeys(cmd.(*CmdKeys))
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
	case cmdLLen:
		return this.handleCommandLLen(cmd.(*CmdLLen))
	case cmdExpire:
		return this.handleCommandExpire(cmd.(*CmdExpire))
	case cmdTTL:
		return this.handleCommandTTL(cmd.(*CmdTTL))
	case cmdHSet:
		return this.handleCommandHSet(cmd.(*CmdHSet))
	case cmdHGet:
		return this.handleCommandHGet(cmd.(*CmdHGet))
	case cmdHDel:
		return this.handleCommandHDel(cmd.(*CmdHDel))
	case cmdHKeys:
		return this.handleCommandHKeys(cmd.(*CmdHKeys))
	case cmdHGetAll:
		return this.handleCommandHGetAll(cmd.(*CmdHGetAll))
	default:
		this.actorCtx.Log().Warning("Received unsupported command [cmd-id: %d]", cmd.CmdId())
		return nil, errors.New("UNSUPPORTED_COMMAND")
	}
}

func (this *storage) ReportStats() {
	ttlSize := this.ttl.Size()
	kvSize := this.kv.Size()

	this.actorCtx.Metrics().ReportStorageTTLSize(ttlSize)
	this.actorCtx.Metrics().ReportStorageKVSize(kvSize)
}

func (this *storage) PurgeTimedOut() {
	purgedCount := 0
	start := time.Now()
	for {
		key, exists := this.ttl.FetchTimedOut(this.tickIdx)
		if !exists {
			break
		}

		//this.actorCtx.Log().Debug("purging timed out key '%s'", key)
		this.kv.Del(key)
		purgedCount++
	}
	if purgedCount > 0 {
		elapsed := time.Since(start)
		//this.actorCtx.Log().Debug("PURGE [bucket: %d; count: %d; elapsed: %v]", this, purgedCount, elapsed)
		this.actorCtx.Metrics().ReportStorageCleanupRecordsCount(purgedCount)
		this.actorCtx.Metrics().ReportStorageCleanupDuration(elapsed)
	}
}

func (this *storage) handleCommandGet(cmd *CmdGet) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)
	if !found {
		return respvalues.NewNil(), nil
	}

	validThru := kve.ValidThru()
	if validThru != ValidThruInfinity && validThru < this.tickIdx {
		this.kv.Del(cmd.key)
		this.ttl.SetTTL(cmd.key, ValidThruInfinity)
		return respvalues.NewNil(), nil
	}

	return kve.Value().ToRESP(), nil
}

func (this *storage) handleCommandSet(cmd *CmdSet) (respvalues.RESPValue, error) {
	this.kv.Set(cmd.key, data.NewScalar(cmd.value), ValidThruInfinity)
	this.ttl.SetTTL(cmd.key, ValidThruInfinity)

	return respvalues.NewStr("OK"), nil
}

func (this *storage) handleCommandDel(cmd *CmdDel) (respvalues.RESPValue, error) {
	kve, existed := this.kv.Get(cmd.key)

	affectedRecords := int64(0)

	if existed {
		this.kv.Del(cmd.key)
		this.ttl.SetTTL(cmd.key, int64(ValidThruInfinity))

		validThru := kve.ValidThru()
		if validThru == ValidThruInfinity || existed && validThru >= this.tickIdx {
			affectedRecords = 1
		}
	}

	return respvalues.NewInt(affectedRecords), nil
}

func (this *storage) handleCommandLLen(cmd *CmdLLen) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if found {
		switch kve.Value().(type) {
		case (*data.ListValue):
			l := kve.Value().(*data.ListValue)
			return respvalues.NewInt(int64(l.Len())), nil
		default:
			return respvalues.NewErr("LLEN: incompatible existing value for this operation"), nil
		}
	} else {
		return respvalues.NewInt(0), nil
	}
}

func (this *storage) handleCommandLPushBack(cmd *CmdLPushBack) (respvalues.RESPValue, error) {
	return this.handleCommandLPushCommon(cmd.key, cmd.value, false)
}

func (this *storage) handleCommandLPushFront(cmd *CmdLPushFront) (respvalues.RESPValue, error) {
	return this.handleCommandLPushCommon(cmd.key, cmd.value, true)
}

func (this *storage) handleCommandLPushCommon(key string, value []byte, front bool) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(key)
	if found {
		switch kve.Value().(type) {
		case (*data.ListValue):
			l := kve.Value().(*data.ListValue)
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
		this.kv.Set(key, l, ValidThruInfinity)

		return respvalues.NewInt(int64(1)), nil
	}
}

func (this *storage) handleCommandLPopBack(cmd *CmdLPopBack) (respvalues.RESPValue, error) {
	return this.handleCommandLPopCommon(cmd.key, false)
}

func (this *storage) handleCommandLPopFront(cmd *CmdLPopFront) (respvalues.RESPValue, error) {
	return this.handleCommandLPopCommon(cmd.key, false)
}

func (this *storage) handleCommandLPopCommon(key string, front bool) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(key)
	if !found {
		return respvalues.NewNil(), nil
	}

	switch kve.Value().(type) {
	case (*data.ListValue):
		l := kve.Value().(*data.ListValue)
		var value []byte = nil
		isEmpty := false
		if front {
			value, isEmpty = l.PopFront()
		} else {
			value, isEmpty = l.PopBack()
		}
		if isEmpty {
			this.kv.Del(key)
			this.ttl.SetTTL(key, ValidThruInfinity)
		}
		return respvalues.NewBulkStr(value), nil
	default:
		return respvalues.NewErr("LPOP*: incompatible existing value for this operation"), nil
	}
}

func (this *storage) handleCommandLGetNth(cmd *CmdLGetNth) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)
	if !found {
		return respvalues.NewNil(), nil
	}

	switch kve.Value().(type) {
	case (*data.ListValue):
		l := kve.Value().(*data.ListValue)
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

func (this *storage) handleCommandTTL(cmd *CmdTTL) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if !found {
		return respvalues.NewInt(-2), nil
	}
	validThru := kve.ValidThru()

	if validThru == ValidThruInfinity {
		return respvalues.NewInt(ValidThruInfinity), nil
	}

	nanosLeft := (validThru - this.tickIdx) * metronome.TickDurationNanos

	if nanosLeft < 0 {
		return respvalues.NewInt(-2), nil
	}
	if cmd.useSeconds {
		return respvalues.NewInt(int64(nanosLeft / 1000000000)), nil
	} else {
		return respvalues.NewInt(int64(nanosLeft / 1000000)), nil
	}
}

func (this *storage) handleCommandExpire(cmd *CmdExpire) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if !found {
		return respvalues.NewInt(0), nil
	}

	validThru := int64(ValidThruInfinity)
	if cmd.expiry != ValidThruInfinity {
		expiryTicks := cmd.expiry * 1000000 / metronome.TickDurationNanos
		validThru = this.tickIdx + expiryTicks
	}

	this.kv.Set(cmd.key, kve.Value(), validThru)
	this.ttl.SetTTL(cmd.key, validThru)

	return respvalues.NewInt(1), nil
}

func (this *storage) handleCommandHSet(cmd *CmdHSet) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if !found {
		dict := data.NewDict()
		dict.Set(cmd.hkey, cmd.hvalue)
		this.kv.Set(cmd.key, dict, ValidThruInfinity)
		return respvalues.NewInt(1), nil
	}

	switch kve.Value().(type) {
	case (*data.DictValue):
		dict := kve.Value().(*data.DictValue)
		keyCreated := dict.Set(cmd.hkey, cmd.hvalue)
		if keyCreated {
			return respvalues.NewInt(1), nil
		} else {
			return respvalues.NewInt(0), nil
		}
	default:
		return respvalues.NewErr("HSET: incompatible existing value for this operation"), nil
	}
}

func (this *storage) handleCommandHGet(cmd *CmdHGet) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if !found {
		return respvalues.NewNil(), nil
	}

	switch kve.Value().(type) {
	case (*data.DictValue):
		dict := kve.Value().(*data.DictValue)
		hvalue, hfound := dict.Get(cmd.hkey)
		if !hfound {
			return respvalues.NewNil(), nil
		} else {
			return respvalues.NewBulkStr(hvalue), nil
		}

	default:
		return respvalues.NewErr("HGET: incompatible existing value for this operation"), nil
	}
}

func (this *storage) handleCommandHDel(cmd *CmdHDel) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if !found {
		return respvalues.NewInt(0), nil
	}

	switch kve.Value().(type) {
	case (*data.DictValue):
		dict := kve.Value().(*data.DictValue)
		hexisted, hempty := dict.Del(cmd.hkey)

		if hempty {
			this.kv.Del(cmd.key)
		}

		if !hexisted {
			return respvalues.NewInt(0), nil
		} else {
			return respvalues.NewInt(1), nil
		}

	default:
		return respvalues.NewErr("HDEL: incompatible existing value for this operation"), nil
	}
}

func (this *storage) handleCommandHKeys(cmd *CmdHKeys) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if !found {
		return respvalues.NewNil(), nil
	}

	switch kve.Value().(type) {
	case (*data.DictValue):
		dict := kve.Value().(*data.DictValue)
		hkeys := dict.Keys()
		hkeysAsResp := list.New()
		for i := 0; i < len(hkeys); i++ {
			k := respvalues.NewStr(hkeys[i])
			hkeysAsResp.PushBack(k)
		}
		return respvalues.NewArray(hkeysAsResp), nil

	default:
		return respvalues.NewErr("HKEYS: incompatible existing value for this operation"), nil
	}
}

func (this *storage) handleCommandHGetAll(cmd *CmdHGetAll) (respvalues.RESPValue, error) {
	kve, found := this.kv.Get(cmd.key)

	if !found {
		return respvalues.NewNil(), nil
	}

	switch kve.Value().(type) {
	case (*data.DictValue):
		dict := kve.Value().(*data.DictValue)
		hvalues := dict.Values()
		hvaluesAsResp := list.New()
		for i := 0; i < len(hvalues); i++ {
			v := respvalues.NewBulkStr(hvalues[i])
			hvaluesAsResp.PushBack(v)
		}
		return respvalues.NewArray(hvaluesAsResp), nil

	default:
		return respvalues.NewErr("HGETALL: incompatible existing value for this operation"), nil
	}
}

func (this *storage) handleCommandKeys(cmd *CmdKeys) (respvalues.RESPValue, error) {
	allKeys := this.kv.Keys()

	fut := respvalues.NewFuture()

	go processKeys(cmd.pattern, allKeys, fut)

	return fut, nil
}

func processKeys(pattern string, allKeys *gtreap.Treap, fut *respvalues.RESPFuture) {
	firstKey := allKeys.Min()
	resultList := list.New()

	if firstKey != nil {
		re, err := regexp.Compile(pattern)
		if err != nil {
			fut.Chan() <- respvalues.NewErr(err.Error())
			return
		}
		allKeys.VisitAscend(firstKey, func(item gtreap.Item) bool {
			key := item.(string)
			if re.Match([]byte(key)) {
				resultList.PushBack(respvalues.NewStr(key))
			}
			return true
		})
	}
	fut.Chan() <- respvalues.NewArray(resultList)
}
