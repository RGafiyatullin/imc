package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type LPopFHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *LPopFHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandLPopFDuration(elapsed)
}

func (this *LPopFHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("LPOPF: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX
	key := reqElements[1].(*respvalues.BasicBulkStr)
	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdLPopFront(key.String()))

	return result
}

func (this *LPopFHandler) Register(registry map[string]CommandHandler) {
	registry["LPOPF"] = this
}

func NewLPopFHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(LPopFHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
