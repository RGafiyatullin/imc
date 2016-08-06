package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type SetHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *SetHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandSetDuration(elapsed)
}

func (this *SetHandler) Handle(req *respvalues.RESPArray) respvalues.RESPValue {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) < 3 {
		return respvalues.NewErr("SET: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	// XXX
	key := reqElements[1].(*respvalues.RESPBulkStr)
	// XXX
	value := reqElements[2].(*respvalues.RESPBulkStr)

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]

	result := bucketApi.RunCmd(bucket.NewCmdSet(key.String(), value.Bytes()))

	return result
}

func (this *SetHandler) Register(registry map[string]CommandHandler) {
	registry["SET"] = this
}

func NewSetHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(SetHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
