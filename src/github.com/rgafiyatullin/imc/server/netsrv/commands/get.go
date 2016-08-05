package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
)

type GetHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *GetHandler) Handle(req *types.BasicArr) types.BasicType {
	return types.NewErr("GET: not implemented")
}

func (this *GetHandler) Register(registry map[string]CommandHandler) {
	registry["GET"] = this
}

func NewGetHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(GetHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
