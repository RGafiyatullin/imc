package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type HKeysHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *HKeysHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandHKeysDuration(elapsed)
}

func (this *HKeysHandler) Handle(req *respvalues.RESPArray) respvalues.RESPValue {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("HKEYS: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX
	key := reqElements[1].(*respvalues.RESPBulkStr)
	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdHKeys(key.String()))

	return result
}

func (this *HKeysHandler) Register(registry map[string]CommandHandler) {
	registry["HKEYS"] = this
}

func NewHKeysHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(HKeysHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
