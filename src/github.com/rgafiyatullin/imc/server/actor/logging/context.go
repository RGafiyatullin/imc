package logging

import "fmt"

type Ctx interface {
	CloneWithName(name string) Ctx
	Debug(fmtStr string, args ...interface{})
	Info(fmtStr string, args ...interface{})
	Warn(fmtStr string, args ...interface{})
	Error(fmtStr string, args ...interface{})
}

type StdoutCtx struct {
	name_ string
}

func (this *StdoutCtx) CloneWithName(name string) Ctx {
	clone := new(StdoutCtx)
	clone.name_ = name
	return clone
}

func (this *StdoutCtx) Debug(fmtStr string, args ...interface{}) {
	fmt.Printf("[%s] Log: [debug] fmt: %v; args: %v\n", this.name_, fmtStr, args)
}

func (this *StdoutCtx) Info(fmtStr string, args ...interface{}) {
	fmt.Printf("[%s] Log: [info] fmt: %v; args: %v\n", this.name_, fmtStr, args)
}

func (this *StdoutCtx) Warn(fmtStr string, args ...interface{}) {
	fmt.Printf("[%s] Log: [warn] fmt: %v; args: %v\n", this.name_, fmtStr, args)
}

func (this *StdoutCtx) Error(fmtStr string, args ...interface{}) {
	fmt.Printf("[%s] Log: [error] fmt: %v; args: %v\n", this.name_, fmtStr, args)
}

func NewStdoutCtx() Ctx {
	stdoutCtx := new(StdoutCtx)
	return stdoutCtx
}
