package logging

import (
	"github.com/rgafiyatullin/imc/server/config"
	"time"
)

// Logging context.
// Proxies the logging requests additionally injecting current actor info into the log message.
type Ctx interface {
	// Creates another context with different actor-name
	CloneWithName(name string) Ctx
	// Blocks until the previously (in terms of the current execution context) emitted log-messages are actually logged.
	// May be useful in case of logging fatal messages prior to requesting system halt.
	Flush()

	Debug(fmtStr string, args ...interface{})
	Info(fmtStr string, args ...interface{})
	Warning(fmtStr string, args ...interface{})
	Error(fmtStr string, args ...interface{})
	Fatal(fmtStr string, args ...interface{})
}

type ctx struct {
	handler_ Handler
	name_    string
}

func (this *ctx) CloneWithName(name string) Ctx {
	clone := new(ctx)
	clone.handler_ = this.handler_
	clone.name_ = name
	return clone
}

func (this *ctx) message(lvl int, fmtStr string, args []interface{}) {
	now := time.Now()
	report := new(LogReport)
	report.level = lvl
	report.at = now
	report.fmt = fmtStr
	report.args = args
	if this.name_ == "" {
		report.entity = "/"
	} else {
		report.entity = this.name_
	}

	this.handler_.Report(report)
}

func (this *ctx) Debug(fmtStr string, args ...interface{}) {
	this.message(lvlDebug, fmtStr, args)
}

func (this *ctx) Info(fmtStr string, args ...interface{}) {
	this.message(lvlInfo, fmtStr, args)
}

func (this *ctx) Warning(fmtStr string, args ...interface{}) {
	this.message(lvlWarning, fmtStr, args)
}

func (this *ctx) Error(fmtStr string, args ...interface{}) {
	this.message(lvlError, fmtStr, args)
}

func (this *ctx) Fatal(fmtStr string, args ...interface{}) {
	this.message(lvlFatal, fmtStr, args)
	this.Flush()
}

func (this *ctx) Flush() {
	this.handler_.Flush()
}

func New(config config.Config) Ctx {
	handler := NewHandler()
	stdoutCtx := new(ctx)
	stdoutCtx.handler_ = handler
	return stdoutCtx
}
