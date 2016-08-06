package logging

import (
	"fmt"
	"time"
)

const ChanBufSize = 64

type LogRequest interface {}

type LogReport struct {
	at     time.Time
	level  int
	fmt    string
	args   []interface{}
	entity string
}

type FlushMsg struct {
	ackTo chan<-bool
}

type Handler interface {
	Report(report *LogReport)
	Flush()
}
type handler struct {
	ch chan<- LogRequest
}

func (this *handler) Report(report *LogReport) {
	this.ch <- report
}

func (this *handler) Flush() {
	ch := make(chan bool, 1)
	req := new(FlushMsg)
	req.ackTo = ch
	this.ch <- req
	<- ch
}

func NewHandler() Handler {
	ch := make(chan LogRequest, ChanBufSize)
	h := new(handler)
	h.ch = ch

	go hanlerEnterLoop(ch)

	return h
}

type logHandlerState struct {
	reports <-chan LogRequest
}

func (this *logHandlerState) init() {
	fmt.Printf("[%v] initializing\n", time.Now().UTC())
}
func (this *logHandlerState) loop() {
	for {
		request := <-this.reports
		switch request.(type) {
		case (*LogReport):
			this.processReport(request.(*LogReport))
		case (*FlushMsg):
			this.processFlush()
			request.(*FlushMsg).ackTo <- true
		}

	}
}
func (this *logHandlerState) processReport(report *LogReport) {
	fmtStr := "[%-39s] [%s] [%s] " + report.fmt + "\n"
	args := make([]interface{}, len(report.args)+3)
	args[0] = report.at.UTC()
	args[1] = levelToString(report.level)
	args[2] = report.entity
	copy(args[3:], report.args)
	fmt.Printf(fmtStr, args...)
}

func (this *logHandlerState) processFlush() {

}

func hanlerEnterLoop(reports chan LogRequest) {
	state := new(logHandlerState)
	state.reports = reports
	state.init()
	state.loop()
}

func levelToString(l int) string {
	switch l {
	case lvlTrace:
		return "TRACE"
	case lvlDebug:
		return "DEBUG"
	case lvlInfo:
		return "INFO_"
	case lvlWarning:
		return "WARN_"
	case lvlError:
		return "ERROR"
	case lvlFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("UNKNOWN:%d", l)
	}
}
