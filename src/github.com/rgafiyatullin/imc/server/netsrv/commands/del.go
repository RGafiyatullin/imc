package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
)

type DelHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *DelHandler) Handle(req *types.BasicArr) types.BasicType {
	reqElements := req.Elements()

	if len(reqElements) < 2 {
		return types.NewErr("DEL: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	affectedRecords := types.NewInt(int64(0))

	for i := 1; i < len(reqElements); i ++ {
		// XXX: sorry
		key := reqElements[i].(*types.BasicBulkStr)
		keyHash := ringmgr.CalcKeyHash(key)
		bucketIdx := keyHash % uint32(len(buckets))
		bucketApi := buckets[bucketIdx]
		keyResult := bucketApi.RunCmd(bucket.NewCmdDel(key.String()))
		affectedRecords = affectedRecords.Add(keyResult.(*types.BasicInt))
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
