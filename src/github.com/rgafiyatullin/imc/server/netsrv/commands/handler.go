package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
)

// A common interface for the client to server commands.
type CommandHandler interface {
	Register(map[string]CommandHandler)
	Handle(req *respvalues.RESPArray) respvalues.RESPValue
}
