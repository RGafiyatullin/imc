package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type HSetHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *HSetHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandHSetDuration(elapsed)
}

func (this *HSetHandler) Handle(req *respvalues.RESPArray) respvalues.RESPValue {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 4 {
		return respvalues.NewErr("HSET: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	// XXX
	key := reqElements[1].(*respvalues.RESPBulkStr)
	hkey := reqElements[2].(*respvalues.RESPBulkStr)
	hvalue := reqElements[3].(*respvalues.RESPBulkStr)

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]

	result := bucketApi.RunCmd(bucket.NewCmdHSet(key.String(), hkey.String(), hvalue.Bytes()))

	return result
}

func (this *HSetHandler) Register(registry map[string]CommandHandler) {
	registry["HSET"] = this
}

func NewHSetHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(HSetHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
