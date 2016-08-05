package commands

import "github.com/rgafiyatullin/imc/protocol/resp/types"

type PingHandler struct {}

func (this *PingHandler) Handle(req *types.BasicArr) types.BasicType {
	return types.NewStr("PONG")
}

func (this *PingHandler) Register(registry map[string]CommandHandler) {
	registry["PING"] = this
}

func NewPingHandler() CommandHandler {
	return new(PingHandler)
}
