package logging

import (
	"github.com/rgafiyatullin/imc/server/config"
	"time"
)

type Ctx interface {
	CloneWithName(name string) Ctx
	Debug(fmtStr string, args ...interface{})
	Info(fmtStr string, args ...interface{})
	Warning(fmtStr string, args ...interface{})
	Error(fmtStr string, args ...interface{})
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

func New(config config.Config) Ctx {
	handler := NewHandler()
	stdoutCtx := new(ctx)
	stdoutCtx.handler_ = handler
	return stdoutCtx
}
