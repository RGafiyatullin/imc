package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
)

type CommandHandler interface {
	Register(map[string]CommandHandler)
	Handle(req *respvalues.BasicArr) respvalues.BasicType
}
