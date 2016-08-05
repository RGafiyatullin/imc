package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"strconv"
	"time"
)

type LNthHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *LNthHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandLGetNthDuration(elapsed)
}

func (this *LNthHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 3 {
		return respvalues.NewErr("LNTH: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX: sorry
	key := reqElements[1].(*respvalues.BasicBulkStr)
	// XXX: sorry
	idx, idxParseErr := strconv.ParseInt(reqElements[2].(*respvalues.BasicBulkStr).String(), 10, 32)
	if idxParseErr != nil {
		return respvalues.NewErr("LNTH: invalid idx specified")
	}

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdLGetNth(key.String(), int(idx)))

	return result
}

func (this *LNthHandler) Register(registry map[string]CommandHandler) {
	registry["LNTH"] = this
}

func NewLNthHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(LNthHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
