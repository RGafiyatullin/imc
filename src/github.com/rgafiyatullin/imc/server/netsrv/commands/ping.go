package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
)

type PingHandler struct {
	ctx actor.Ctx
}

func (this *PingHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	return respvalues.NewStr("PONG")
}

func (this *PingHandler) Register(registry map[string]CommandHandler) {
	registry["PING"] = this
}

func NewPingHandler(ctx actor.Ctx) CommandHandler {
	h := new(PingHandler)
	h.ctx = ctx
	return h
}
