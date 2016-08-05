package metrics

import (
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/rcrowley/go-metrics"
	"github.com/rgafiyatullin/imc/server/actor/logging"
	"github.com/rgafiyatullin/imc/server/config"
	"time"
)

type Ctx interface {
	ReportCommandDuration(d time.Duration)
	ReportCommandGetDuration(d time.Duration)
	ReportCommandSetDuration(d time.Duration)
	ReportCommandDelDuration(d time.Duration)
}

type ctx struct {
	log    logging.Ctx
	config config.Config

	command_duration_h metrics.Histogram
	command_rate_m     metrics.Meter
	command_get_duration_h metrics.Histogram
	command_get_rate_m     metrics.Meter
	command_set_duration_h metrics.Histogram
	command_set_rate_m     metrics.Meter
	command_del_duration_h metrics.Histogram
	command_del_rate_m     metrics.Meter
}

func (this *ctx) ReportCommandDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_rate_m.Mark(1)
	this.command_duration_h.Update(us)
}

func (this *ctx) ReportCommandGetDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_get_rate_m.Mark(1)
	this.command_get_duration_h.Update(us)
}

func (this *ctx) ReportCommandSetDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_set_rate_m.Mark(1)
	this.command_set_duration_h.Update(us)
}

func (this *ctx) ReportCommandDelDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_del_rate_m.Mark(1)
	this.command_del_duration_h.Update(us)
}

func (this *ctx) init(log logging.Ctx, config config.Config) {
	this.log = log
	this.config = config

	this.log.Info("init")

	sample := metrics.NewExpDecaySample(1028, 0.015)

	this.command_duration_h = metrics.NewHistogram(sample)
	this.command_rate_m = metrics.NewMeter()

	this.command_get_duration_h = metrics.NewHistogram(sample)
	this.command_get_rate_m = metrics.NewMeter()

	this.command_set_duration_h = metrics.NewHistogram(sample)
	this.command_set_rate_m = metrics.NewMeter()

	this.command_del_duration_h = metrics.NewHistogram(sample)
	this.command_del_rate_m = metrics.NewMeter()


	metrics.DefaultRegistry.Register("netsrv.command.all.duration.h", this.command_duration_h)
	metrics.DefaultRegistry.Register("netsrv.command.rate.m", this.command_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.get.duration.h", this.command_get_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.get.rate.m", this.command_get_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.set.duration.h", this.command_set_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.set.rate.m", this.command_set_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.del.duration.h", this.command_del_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.del.rate.m", this.command_del_rate_m)
}

func (this *ctx) startGraphiteReporter() {
	if this.config.Metrics().GraphiteEnabled() {
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
