package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type LPshBHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *LPshBHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandLPshBDuration(elapsed)
}

func (this *LPshBHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) < 3 {
		return respvalues.NewErr("LPSHB: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	// XXX: sorry
	key := reqElements[1].(*respvalues.BasicBulkStr)
	// XXX: sorry again
	value := reqElements[2].(*respvalues.BasicBulkStr)

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]

	result := bucketApi.RunCmd(bucket.NewCmdLPushBack(key.String(), value.Bytes()))

	return result
}

func (this *LPshBHandler) Register(registry map[string]CommandHandler) {
	registry["LPSHB"] = this
}

func NewLPshBHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(LPshBHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
