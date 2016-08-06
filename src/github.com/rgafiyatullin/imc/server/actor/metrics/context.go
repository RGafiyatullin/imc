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
	ReportCommandLPshFDuration(d time.Duration)
	ReportCommandLPshBDuration(d time.Duration)
	ReportCommandLPopFDuration(d time.Duration)
	ReportCommandLPopBDuration(d time.Duration)
	ReportCommandLGetNthDuration(d time.Duration)
	ReportCommandExpireDuration(d time.Duration)
	ReportCommandTTLDuration(d time.Duration)
	ReportCommandHSetDuration(d time.Duration)
	ReportCommandHGetDuration(d time.Duration)
	ReportCommandHDelDuration(d time.Duration)
	ReportCommandHKeysDuration(d time.Duration)
	ReportCommandHGetAllDuration(d time.Duration)
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

	command_lpshf_duration_h   metrics.Histogram
	command_lpshf_rate_m       metrics.Meter
	command_lpshb_duration_h   metrics.Histogram
	command_lpshb_rate_m       metrics.Meter
	command_lpopf_duration_h   metrics.Histogram
	command_lpopf_rate_m       metrics.Meter
	command_lpopb_duration_h   metrics.Histogram
	command_lpopb_rate_m       metrics.Meter
	command_lgetnth_duration_h metrics.Histogram
	command_lgetnth_rate_m     metrics.Meter

	command_expire_duration_h metrics.Histogram
	command_expire_rate_m     metrics.Meter
	command_ttl_duration_h    metrics.Histogram
	command_ttl_rate_m        metrics.Meter

	command_hset_duration_h    metrics.Histogram
	command_hset_rate_m        metrics.Meter
	command_hget_duration_h    metrics.Histogram
	command_hget_rate_m        metrics.Meter
	command_hdel_duration_h    metrics.Histogram
	command_hdel_rate_m        metrics.Meter
	command_hkeys_duration_h   metrics.Histogram
	command_hkeys_rate_m       metrics.Meter
	command_hgetall_duration_h metrics.Histogram
	command_hgetall_rate_m     metrics.Meter
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

func (this *ctx) ReportCommandLGetNthDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_lgetnth_rate_m.Mark(1)
	this.command_lgetnth_duration_h.Update(us)
}

func (this *ctx) ReportCommandLPopBDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_lpopb_rate_m.Mark(1)
	this.command_lpopb_duration_h.Update(us)
}

func (this *ctx) ReportCommandLPshBDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_lpshb_rate_m.Mark(1)
	this.command_lpshb_duration_h.Update(us)
}

func (this *ctx) ReportCommandLPopFDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_lpopf_rate_m.Mark(1)
	this.command_lpopf_duration_h.Update(us)
}

func (this *ctx) ReportCommandLPshFDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_lpshf_rate_m.Mark(1)
	this.command_lpshf_duration_h.Update(us)
}

func (this *ctx) ReportCommandExpireDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_expire_rate_m.Mark(1)
	this.command_expire_duration_h.Update(us)
}

func (this *ctx) ReportCommandTTLDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_ttl_rate_m.Mark(1)
	this.command_ttl_duration_h.Update(us)
}

func (this *ctx) ReportCommandHSetDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_hset_rate_m.Mark(1)
	this.command_hset_duration_h.Update(us)
}

func (this *ctx) ReportCommandHGetDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_hget_rate_m.Mark(1)
	this.command_hget_duration_h.Update(us)
}

func (this *ctx) ReportCommandHDelDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_hdel_rate_m.Mark(1)
	this.command_hdel_duration_h.Update(us)
}

func (this *ctx) ReportCommandHKeysDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_hkeys_rate_m.Mark(1)
	this.command_hkeys_duration_h.Update(us)
}
func (this *ctx) ReportCommandHGetAllDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_hgetall_rate_m.Mark(1)
	this.command_hgetall_duration_h.Update(us)
}

func (this *ctx) init(log logging.Ctx, config config.Config) {
	this.log = log
	this.config = config

	this.log.Info("init")

	sampleSize := 1028
	alpha := 0.015

	this.command_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_rate_m = metrics.NewMeter()

	this.command_get_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_get_rate_m = metrics.NewMeter()

	this.command_set_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_set_rate_m = metrics.NewMeter()

	this.command_del_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_del_rate_m = metrics.NewMeter()

	this.command_lpshf_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_lpshf_rate_m = metrics.NewMeter()

	this.command_lpshb_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_lpshb_rate_m = metrics.NewMeter()

	this.command_lpopf_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_lpopf_rate_m = metrics.NewMeter()

	this.command_lpopb_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_lpopb_rate_m = metrics.NewMeter()

	this.command_lgetnth_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_lgetnth_rate_m = metrics.NewMeter()

	this.command_expire_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_expire_rate_m = metrics.NewMeter()

	this.command_ttl_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_ttl_rate_m = metrics.NewMeter()

	this.command_hset_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_hset_rate_m = metrics.NewMeter()

	this.command_hget_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_hget_rate_m = metrics.NewMeter()

	this.command_hdel_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_hdel_rate_m = metrics.NewMeter()

	this.command_hkeys_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_hkeys_rate_m = metrics.NewMeter()

	this.command_hgetall_duration_h = metrics.NewHistogram(
		metrics.NewExpDecaySample(sampleSize, alpha))
	this.command_hgetall_rate_m = metrics.NewMeter()

	metrics.DefaultRegistry.Register("netsrv.command.duration.h", this.command_duration_h)
	metrics.DefaultRegistry.Register("netsrv.command.rate.m", this.command_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.get.duration.h", this.command_get_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.get.rate.m", this.command_get_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.set.duration.h", this.command_set_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.set.rate.m", this.command_set_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.del.duration.h", this.command_del_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.del.rate.m", this.command_del_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.lpshf.duration.h", this.command_lpshf_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.lpshf.rate.m", this.command_lpshf_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.lpshb.duration.h", this.command_lpshb_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.lpshb.rate.m", this.command_lpshb_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.lpopf.duration.h", this.command_lpopf_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.lpopf.rate.m", this.command_lpopf_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.lpopb.duration.h", this.command_lpopb_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.lpopb.rate.m", this.command_lpopb_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.lgetnth.duration.h", this.command_lgetnth_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.lgetnth.rate.m", this.command_lgetnth_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.expire.duration.h", this.command_expire_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.expire.rate.m", this.command_expire_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.ttl.duration.h", this.command_ttl_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.ttl.rate.m", this.command_ttl_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.hset.duration.h", this.command_ttl_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.hset.rate.m", this.command_ttl_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.hget.duration.h", this.command_ttl_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.hget.rate.m", this.command_ttl_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.hdel.duration.h", this.command_ttl_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.hdel.rate.m", this.command_ttl_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.hkeys.duration.h", this.command_hkeys_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.hkeys.rate.m", this.command_hkeys_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.hgetall.duration.h", this.command_hgetall_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.hgetall.rate.m", this.command_hgetall_rate_m)
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
