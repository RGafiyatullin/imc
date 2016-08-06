package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"time"
)

type AuthHandler struct {
	ctx      actor.Ctx
	password string
	handlers Handlers
}

func (this *AuthHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandAuthDuration(elapsed)
}

func (this *AuthHandler) Handle(req *respvalues.RESPArray) respvalues.RESPValue {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("AUTH: malformed command")
	}

	password := reqElements[1].(*respvalues.RESPBulkStr).String()

	if password != this.password {
		this.ctx.Metrics().ReportCommandAuthFailure()
		return respvalues.NewInt(0)
	} else {
		this.ctx.Metrics().ReportCommandAuthSuccess()
		this.handlers.InitCommandsFullSet()
		return respvalues.NewInt(1)
	}
}

func (this *AuthHandler) Register(registry map[string]CommandHandler) {
	registry["AUTH"] = this
}

func NewAuthHandler(ctx actor.Ctx, password string, handlers Handlers) CommandHandler {
	h := new(AuthHandler)
	h.ctx = ctx
	h.password = password
	h.handlers = handlers
	return h
}
