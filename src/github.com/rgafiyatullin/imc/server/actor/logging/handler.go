package logging

import (
	"fmt"
	"time"
)

const ChanBufSize = 64

type LogReport struct {
	at     time.Time
	level  int
	fmt    string
	args   []interface{}
	entity string
}

type Handler interface {
	Report(report *LogReport)
}
type handler struct {
	ch chan<- *LogReport
}

func (this *handler) Report(report *LogReport) {
	this.ch <- report
}

func NewHandler() Handler {
	ch := make(chan *LogReport, ChanBufSize)
	h := new(handler)
	h.ch = ch

	go hanlerEnterLoop(ch)

	return h
}

type logHandlerState struct {
	reports <-chan *LogReport
}

func (this *logHandlerState) init() {
	fmt.Printf("[%v] initializing\n", time.Now().UTC())
}
func (this *logHandlerState) loop() {
	for {
		report := <-this.reports
		this.processReport(report)
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

func hanlerEnterLoop(reports chan *LogReport) {
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
