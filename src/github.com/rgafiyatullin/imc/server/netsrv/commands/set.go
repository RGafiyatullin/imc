package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
)

type SetHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *SetHandler) Handle(req *types.BasicArr) types.BasicType {
	return types.NewErr("SET: not implemented")
}

func (this *SetHandler) Register(registry map[string]CommandHandler) {
	registry["SET"] = this
}

func NewSetHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(SetHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
