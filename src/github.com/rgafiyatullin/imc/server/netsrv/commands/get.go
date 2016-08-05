package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
)

type GetHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *GetHandler) Handle(req *types.BasicArr) types.BasicType {
	reqElements := req.Elements()

	if len(reqElements) != 2 {
		return types.NewErr("GET: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX: sorry
	key := reqElements[1].(*types.BasicBulkStr)
	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdGet(key.String()))

	return result
}

func (this *GetHandler) Register(registry map[string]CommandHandler) {
	registry["GET"] = this
}

func NewGetHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(GetHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
