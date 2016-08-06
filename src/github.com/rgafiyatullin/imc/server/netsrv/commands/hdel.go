package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type HDelHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *HDelHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandHDelDuration(elapsed)
}

func (this *HDelHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 3 {
		return respvalues.NewErr("HDEL: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	// XXX
	key := reqElements[1].(*respvalues.BasicBulkStr)
	hkey := reqElements[2].(*respvalues.BasicBulkStr)

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]

	result := bucketApi.RunCmd(bucket.NewCmdHDel(key.String(), hkey.String()))

	return result
}

func (this *HDelHandler) Register(registry map[string]CommandHandler) {
	registry["HDEL"] = this
}

func NewHDelHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(HDelHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
