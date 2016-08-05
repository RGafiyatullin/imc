package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
)

type SetHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *SetHandler) Handle(req *types.BasicArr) types.BasicType {
	reqElements := req.Elements()

	if len(reqElements) < 3 {
		return types.NewErr("SET: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()

	expiry := uint64(0)
	// XXX: sorry
	key := reqElements[1].(*types.BasicBulkStr)
	value := reqElements[2]

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdSet(key.String(), value, expiry))

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
