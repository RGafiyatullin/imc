package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
)

// An interface to switch the predefined commands sets
type Handlers interface {
	InitCommands()
	InitCommandsUnauthed()
	InitCommandsFullSet()
}

// A common interface for the client to server commands.
type CommandHandler interface {
	Register(map[string]CommandHandler)
	Handle(req *respvalues.RESPArray) respvalues.RESPValue
}
