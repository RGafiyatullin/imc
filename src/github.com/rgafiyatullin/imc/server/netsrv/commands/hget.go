package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type HGetHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *HGetHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandHGetDuration(elapsed)
}

func (this *HGetHandler) Handle(req *respvalues.RESPArray) respvalues.RESPValue {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 3 {
		return respvalues.NewErr("HGET: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	// XXX
	key := reqElements[1].(*respvalues.RESPBulkStr)
	hkey := reqElements[2].(*respvalues.RESPBulkStr)

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]

	result := bucketApi.RunCmd(bucket.NewCmdHGet(key.String(), hkey.String()))

	return result
}

func (this *HGetHandler) Register(registry map[string]CommandHandler) {
	registry["HGET"] = this
}

func NewHGetHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(HGetHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
