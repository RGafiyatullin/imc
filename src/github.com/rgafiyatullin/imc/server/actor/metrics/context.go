package metrics

import (
	"github.com/rcrowley/go-metrics"
	"github.com/rgafiyatullin/imc/server/actor/logging"
	"time"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/cyberdelia/go-metrics-graphite"
)

type Ctx interface {
	ReportCommandDuration(d time.Duration)
}

type ctx struct {
	log logging.Ctx
	config config.Config

	command_duration_h metrics.Histogram
	command_rate_m metrics.Meter
}

func (this *ctx) ReportCommandDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_rate_m.Mark(1)
	this.command_duration_h.Update(us)
}

func (this *ctx) init(log logging.Ctx, config config.Config) {
	this.log = log
	this.config = config

	this.log.Info("init")

	sample := metrics.NewExpDecaySample(1028, 0.015)
	this.command_duration_h = metrics.NewHistogram(sample)
	this.command_rate_m = metrics.NewMeter()

	metrics.DefaultRegistry.Register("netsrv.command.duration.h", this.command_duration_h)
	metrics.DefaultRegistry.Register("netsrv.command.rate.m", this.command_rate_m)
}

func (this *ctx) startGraphiteReporter() {
	if (this.config.Metrics().GraphiteEnabled()) {
		addr := this.config.Metrics().GraphiteAddr()
		prefix := this.config.Metrics().GraphitePrefix()

		this.log.Info("starting up graphite reporter [addr: %s; prefix: %s]", addr, prefix)
		go graphite.Graphite(metrics.DefaultRegistry, 10e9, prefix, addr)
	} else {
		this.log.Info("graphite reporter disabled; not starting it")
	}
}

func New(log logging.Ctx, config config.Config) Ctx {
	ctx := new(ctx)
	ctx.init(log, config)
	ctx.startGraphiteReporter()
	return ctx
}

