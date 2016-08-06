package metrics

import (
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/rcrowley/go-metrics"
	"github.com/rgafiyatullin/imc/server/actor/logging"
	"github.com/rgafiyatullin/imc/server/config"
	"runtime"
	"time"
)

// A context to report metrics via
type Ctx interface {
	ReportCommandDuration(d time.Duration)
	ReportCommandGetDuration(d time.Duration)
	ReportCommandSetDuration(d time.Duration)
	ReportCommandDelDuration(d time.Duration)
	ReportCommandKeysDuration(d time.Duration)
	ReportCommandLPshFDuration(d time.Duration)
	ReportCommandLPshBDuration(d time.Duration)
	ReportCommandLPopFDuration(d time.Duration)
	ReportCommandLPopBDuration(d time.Duration)
	ReportCommandLGetNthDuration(d time.Duration)
	ReportCommandLLenDuration(d time.Duration)
	ReportCommandExpireDuration(d time.Duration)
	ReportCommandTTLDuration(d time.Duration)
	ReportCommandHSetDuration(d time.Duration)
	ReportCommandHGetDuration(d time.Duration)
	ReportCommandHDelDuration(d time.Duration)
	ReportCommandHKeysDuration(d time.Duration)
	ReportCommandHGetAllDuration(d time.Duration)
	ReportCommandAuthDuration(d time.Duration)
	ReportCommandAuthSuccess()
	ReportCommandAuthFailure()
	ReportConnUp()
	ReportConnDn()
	ReportConnCount(c int)
	ReportStorageCleanupRecordsCount(c int)
	ReportStorageCleanupDuration(d time.Duration)
	ReportStorageTTLSize(c int)
	ReportStorageKVSize(c int)
}

type ctx struct {
	log    logging.Ctx
	config config.Config

	go_runtime_numcpu_g       metrics.Gauge
	go_runtime_maxprocs_g     metrics.Gauge
	go_runtime_numgoroutine_g metrics.Gauge

	storage_cleanup_count_m    metrics.Meter
	storage_cleanup_count_h    metrics.Histogram
	storage_cleanup_duration_h metrics.Histogram

	storage_bucket_ttl_size_h metrics.Histogram
	storage_bucket_kv_size_h  metrics.Histogram

	conn_count_g metrics.Gauge
	conn_up_m    metrics.Meter
	conn_dn_m    metrics.Meter

	command_duration_h metrics.Histogram
	command_rate_m     metrics.Meter

	command_get_duration_h  metrics.Histogram
	command_get_rate_m      metrics.Meter
	command_set_duration_h  metrics.Histogram
	command_set_rate_m      metrics.Meter
	command_del_duration_h  metrics.Histogram
	command_del_rate_m      metrics.Meter
	command_keys_duration_h metrics.Histogram
	command_keys_rate_m     metrics.Meter

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
	command_llen_rate_m        metrics.Meter
	command_llen_duration_h    metrics.Histogram

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

	command_auth_duration_h metrics.Histogram
	command_auth_rate_m     metrics.Meter
	command_auth_success_m  metrics.Meter
	command_auth_failure_m  metrics.Meter
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

func (this *ctx) ReportCommandKeysDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_keys_rate_m.Mark(1)
	this.command_keys_duration_h.Update(us)
}

func (this *ctx) ReportCommandLGetNthDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_lgetnth_rate_m.Mark(1)
	this.command_lgetnth_duration_h.Update(us)
}

func (this *ctx) ReportCommandLLenDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_llen_rate_m.Mark(1)
	this.command_llen_duration_h.Update(us)
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

func (this *ctx) ReportCommandAuthDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.command_auth_rate_m.Mark(1)
	this.command_auth_duration_h.Update(us)
}

func (this *ctx) ReportCommandAuthFailure() {
	this.command_auth_failure_m.Mark(1)
}
func (this *ctx) ReportCommandAuthSuccess() {
	this.command_auth_success_m.Mark(1)
}

func (this *ctx) ReportConnCount(c int) {
	this.conn_count_g.Update(int64(c))
}

func (this *ctx) ReportConnUp() {
	this.conn_up_m.Mark(1)
}
func (this *ctx) ReportConnDn() {
	this.conn_dn_m.Mark(1)
}

func (this *ctx) ReportStorageCleanupRecordsCount(c int) {
	this.storage_cleanup_count_h.Update(int64(c))
	this.storage_cleanup_count_m.Mark(int64(c))
}
func (this *ctx) ReportStorageCleanupDuration(d time.Duration) {
	us := d.Nanoseconds() / 1000
	this.storage_cleanup_duration_h.Update(int64(us))
}

func (this *ctx) ReportStorageTTLSize(c int) {
	this.storage_bucket_ttl_size_h.Update(int64(c))
}
func (this *ctx) ReportStorageKVSize(c int) {
	this.storage_bucket_kv_size_h.Update(int64(c))
}

func reportRuntimeMetrics(cpus metrics.Gauge, maxprocs metrics.Gauge, goroutines metrics.Gauge) {
	for {
		cpus.Update(int64(runtime.NumCPU()))
		maxprocs.Update(int64(runtime.GOMAXPROCS(0)))
		goroutines.Update(int64(runtime.NumGoroutine()))
		time.Sleep(1 * time.Second)
	}
}

func (this *ctx) init(log logging.Ctx, config config.Config) {
	this.log = log
	this.config = config

	this.log.Info("init")

	sampleSize := 1028
	sampleAlpha := 0.15
	sample := metrics.NewExpDecaySample(sampleSize, sampleAlpha)

	this.go_runtime_numcpu_g = metrics.NewGauge()
	this.go_runtime_maxprocs_g = metrics.NewGauge()
	this.go_runtime_numgoroutine_g = metrics.NewGauge()
	go reportRuntimeMetrics(
		this.go_runtime_numcpu_g, this.go_runtime_maxprocs_g, this.go_runtime_numgoroutine_g)

	this.storage_bucket_kv_size_h = metrics.NewHistogram(
		metrics.NewUniformSample(int(this.config.Storage().RingSize())))
	this.storage_bucket_ttl_size_h = metrics.NewHistogram(
		metrics.NewUniformSample(int(this.config.Storage().RingSize())))

	this.storage_cleanup_duration_h = metrics.NewHistogram(sample)
	this.storage_cleanup_count_h = metrics.NewHistogram(sample)
	this.storage_cleanup_count_m = metrics.NewMeter()

	this.conn_count_g = metrics.NewGauge()
	this.conn_up_m = metrics.NewMeter()
	this.conn_dn_m = metrics.NewMeter()

	this.command_duration_h = metrics.NewHistogram(sample)
	this.command_rate_m = metrics.NewMeter()

	this.command_get_duration_h = metrics.NewHistogram(sample)
	this.command_get_rate_m = metrics.NewMeter()

	this.command_set_duration_h = metrics.NewHistogram(sample)
	this.command_set_rate_m = metrics.NewMeter()

	this.command_del_duration_h = metrics.NewHistogram(sample)
	this.command_del_rate_m = metrics.NewMeter()

	this.command_keys_duration_h = metrics.NewHistogram(sample)
	this.command_keys_rate_m = metrics.NewMeter()

	this.command_lpshf_duration_h = metrics.NewHistogram(sample)
	this.command_lpshf_rate_m = metrics.NewMeter()

	this.command_lpshb_duration_h = metrics.NewHistogram(sample)
	this.command_lpshb_rate_m = metrics.NewMeter()

	this.command_lpopf_duration_h = metrics.NewHistogram(sample)
	this.command_lpopf_rate_m = metrics.NewMeter()

	this.command_lpopb_duration_h = metrics.NewHistogram(sample)
	this.command_lpopb_rate_m = metrics.NewMeter()

	this.command_lgetnth_duration_h = metrics.NewHistogram(sample)
	this.command_lgetnth_rate_m = metrics.NewMeter()

	this.command_llen_duration_h = metrics.NewHistogram(sample)
	this.command_llen_rate_m = metrics.NewMeter()

	this.command_expire_duration_h = metrics.NewHistogram(sample)
	this.command_expire_rate_m = metrics.NewMeter()

	this.command_ttl_duration_h = metrics.NewHistogram(sample)
	this.command_ttl_rate_m = metrics.NewMeter()

	this.command_hset_duration_h = metrics.NewHistogram(sample)
	this.command_hset_rate_m = metrics.NewMeter()

	this.command_hget_duration_h = metrics.NewHistogram(sample)
	this.command_hget_rate_m = metrics.NewMeter()

	this.command_hdel_duration_h = metrics.NewHistogram(sample)
	this.command_hdel_rate_m = metrics.NewMeter()

	this.command_hkeys_duration_h = metrics.NewHistogram(sample)
	this.command_hkeys_rate_m = metrics.NewMeter()

	this.command_hgetall_duration_h = metrics.NewHistogram(sample)
	this.command_hgetall_rate_m = metrics.NewMeter()

	this.command_auth_duration_h = metrics.NewHistogram(sample)
	this.command_auth_rate_m = metrics.NewMeter()
	this.command_auth_success_m = metrics.NewMeter()
	this.command_auth_failure_m = metrics.NewMeter()

	metrics.DefaultRegistry.Register("go.runtime.numcpu.g", this.go_runtime_numcpu_g)
	metrics.DefaultRegistry.Register("go.runtime.maxprocs.g", this.go_runtime_maxprocs_g)
	metrics.DefaultRegistry.Register("go.runtime.numgoroutine.g", this.go_runtime_numgoroutine_g)

	metrics.DefaultRegistry.Register("storage.cleanup.duration.h", this.storage_cleanup_duration_h)
	metrics.DefaultRegistry.Register("storage.cleanup.records_count.h", this.storage_cleanup_count_h)
	metrics.DefaultRegistry.Register("storage.cleanup.records_count.m", this.storage_cleanup_count_m)

	metrics.DefaultRegistry.Register("storage.bucket.kv.size.h", this.storage_bucket_kv_size_h)
	metrics.DefaultRegistry.Register("storage.bucket.ttl.size.h", this.storage_bucket_ttl_size_h)

	metrics.DefaultRegistry.Register("netsrv.conn.count.g", this.conn_count_g)
	metrics.DefaultRegistry.Register("netsrv.conn.up.rate.m", this.conn_up_m)
	metrics.DefaultRegistry.Register("netsrv.conn.dn.rate.m", this.conn_dn_m)

	metrics.DefaultRegistry.Register("netsrv.command.duration.h", this.command_duration_h)
	metrics.DefaultRegistry.Register("netsrv.command.rate.m", this.command_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.get.duration.h", this.command_get_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.get.rate.m", this.command_get_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.set.duration.h", this.command_set_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.set.rate.m", this.command_set_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.del.duration.h", this.command_del_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.del.rate.m", this.command_del_rate_m)

	metrics.DefaultRegistry.Register("netsrv.commands.keys.duration.h", this.command_keys_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.keys.rate.m", this.command_keys_rate_m)

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

	metrics.DefaultRegistry.Register("netsrv.commands.auth.duration.h", this.command_auth_duration_h)
	metrics.DefaultRegistry.Register("netsrv.commands.auth.rate.m", this.command_auth_rate_m)
	metrics.DefaultRegistry.Register("netsrv.commands.auth.success_rate.m", this.command_auth_success_m)
	metrics.DefaultRegistry.Register("netsrv.commands.auth.failure_rate.m", this.command_auth_failure_m)
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
