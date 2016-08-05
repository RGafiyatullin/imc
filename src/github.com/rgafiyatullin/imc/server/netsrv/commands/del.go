package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type DelHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *DelHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandDelDuration(elapsed)
}

func (this *DelHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) < 2 {
		return respvalues.NewErr("DEL: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	affectedRecords := respvalues.NewInt(int64(0))

	for i := 1; i < len(reqElements); i++ {
		// XXX
		key := reqElements[i].(*respvalues.BasicBulkStr)
		keyHash := ringmgr.CalcKeyHash(key)
		bucketIdx := keyHash % uint32(len(buckets))
		bucketApi := buckets[bucketIdx]
		keyResult := bucketApi.RunCmd(bucket.NewCmdDel(key.String()))
		affectedRecords = affectedRecords.Plus(keyResult.(*respvalues.BasicInt))
	}

	return affectedRecords
}

func (this *DelHandler) Register(registry map[string]CommandHandler) {
	registry["DEL"] = this
}

func NewDelHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(DelHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
