package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
)

type LPopBHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *LPopBHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandLPopBDuration(elapsed)
}

func (this *LPopBHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("LPOPB: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX: sorry
	key := reqElements[1].(*respvalues.BasicBulkStr)
	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdLPopBack(key.String()))

	return result
}

func (this *LPopBHandler) Register(registry map[string]CommandHandler) {
	registry["LPOPB"] = this
}

func NewLPopBHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(LPopBHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}

