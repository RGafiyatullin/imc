package commands

import (
	"container/list"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"time"
)

type KeysHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *KeysHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandKeysDuration(elapsed)
}

func (this *KeysHandler) Handle(req *respvalues.RESPArray) respvalues.RESPValue {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return respvalues.NewErr("KEYS: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX
	pattern := reqElements[1].(*respvalues.RESPBulkStr)

	resultElements := list.New()

	for bucketIdx := 0; bucketIdx < len(buckets); bucketIdx++ {
		bucketApi := buckets[bucketIdx]
		bucketResult := bucketApi.RunCmd(bucket.NewCmdKeys(pattern.String()))
		switch bucketResult.(type) {
		case *respvalues.RESPFuture:
			respFut := bucketResult.(*respvalues.RESPFuture)
			respPresent := respFut.Await()
			switch respPresent.(type) {
			case *respvalues.RESPArray:
				respArr := respPresent.(*respvalues.RESPArray)
				elements := respArr.Elements()
				for i := 0; i < len(elements); i++ {
					resultElements.PushBack(elements[i])
				}

			case *respvalues.RESPErr:
				return bucketResult

			default:
				this.ctx.Log().Warning(
					"Unexpected response (inside Future) to CmdKeys from bucket[%d]: %+v", bucketIdx, bucketResult)
				return respvalues.NewErr("Internal server error")
			}

		default:
			this.ctx.Log().Warning(
				"Unexpected response (instead of Future) to CmdKeys from bucket[%d]: %+v", bucketIdx, bucketResult)
			return respvalues.NewErr("Internal server error")
		}
	}

	return respvalues.NewArray(resultElements)
}

func (this *KeysHandler) Register(registry map[string]CommandHandler) {
	registry["KEYS"] = this
}

func NewKeysHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(KeysHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
