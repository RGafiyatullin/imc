package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
)

type DelHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *DelHandler) Handle(req *types.BasicArr) types.BasicType {
	return types.NewErr("DEL: not implemented")
}

func (this *DelHandler) Register(registry map[string]CommandHandler) {
	registry["DEL"] = this
}

func NewDelHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(DelHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
