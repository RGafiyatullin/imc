package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type TTLHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *TTLHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandTTLDuration(elapsed)
}

func (this *TTLHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("TTL/PTTL: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX
	cmd := reqElements[0].(*respvalues.BasicBulkStr).String()
	key := reqElements[1].(*respvalues.BasicBulkStr)
	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdTTL(key.String(), cmd == "TTL"))

	return result
}

func (this *TTLHandler) Register(registry map[string]CommandHandler) {
	registry["TTL"] = this
	registry["PTTL"] = this
}

func NewTTLHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(TTLHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
