package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
)

type PingHandler struct {
	ctx actor.Ctx
}

func (this *PingHandler) Handle(req *types.BasicArr) types.BasicType {
	return types.NewStr("PONG")
}

func (this *PingHandler) Register(registry map[string]CommandHandler) {
	registry["PING"] = this
}

func NewPingHandler(ctx actor.Ctx) CommandHandler {
	h := new(PingHandler)
	h.ctx = ctx
	return h
}
