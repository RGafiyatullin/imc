package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type HGetAllHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *HGetAllHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandHGetAllDuration(elapsed)
}

func (this *HGetAllHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("HGETALL: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX
	key := reqElements[1].(*respvalues.BasicBulkStr)
	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdHGetAll(key.String()))

	return result
}

func (this *HGetAllHandler) Register(registry map[string]CommandHandler) {
	registry["HGETALL"] = this
}

func NewHGetAllHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(HGetAllHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
