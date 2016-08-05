package commands

import "github.com/rgafiyatullin/imc/protocol/resp/types"

type CommandHandler interface {
	Register(map[string]CommandHandler)
	Handle(req *types.BasicArr) types.BasicType
}
