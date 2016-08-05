package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
)

type LPshFHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *LPshFHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandLPshFDuration(elapsed)
}

func (this *LPshFHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) < 3 {
		return respvalues.NewErr("LPSHF: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	// XXX: sorry
	key := reqElements[1].(*respvalues.BasicBulkStr)
	// XXX: sorry again
	value := reqElements[2].(*respvalues.BasicBulkStr)

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]

	cmd := bucket.NewCmdLPushFront(key.String(), value.Bytes())
	result := bucketApi.RunCmd(cmd)

	return result
}

func (this *LPshFHandler) Register(registry map[string]CommandHandler) {
	registry["LPSHF"] = this
}

func NewLPshFHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(LPshFHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
