package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type LLenHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *LLenHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandLLenDuration(elapsed)
}

func (this *LLenHandler) Handle(req *respvalues.RESPArray) respvalues.RESPValue {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("LLEN: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX
	key := reqElements[1].(*respvalues.RESPBulkStr)
	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdLLen(key.String()))

	return result
}

func (this *LLenHandler) Register(registry map[string]CommandHandler) {
	registry["LLEN"] = this
}

func NewLLenHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(LLenHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
