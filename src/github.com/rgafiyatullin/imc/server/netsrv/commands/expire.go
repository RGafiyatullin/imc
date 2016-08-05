package commands

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"strconv"
	"time"
)

type ExpireHandler struct {
	ctx     actor.Ctx
	ringMgr ringmgr.RingMgr
}

func (this *ExpireHandler) reportTime(start time.Time) {
	elapsed := time.Since(start)
	this.ctx.Metrics().ReportCommandExpireDuration(elapsed)
}

func (this *ExpireHandler) Handle(req *respvalues.BasicArr) respvalues.BasicType {
	startTime := time.Now()
	defer this.reportTime(startTime)

	reqElements := req.Elements()
	cmd := reqElements[0].(*respvalues.BasicBulkStr).String()

	if len(reqElements) != 3 && (cmd == "EXPIRE" || cmd == "PEXPIRE") {
		return respvalues.NewErr("EXPIRE/PEXPIRE: malformed command")
	}
	if len(reqElements) != 2 && cmd == "PERSIST" {
		return respvalues.NewErr("PERSIST: malformed command")
	}

	buckets := this.ringMgr.QueryBuckets()
	// XXX: sorry

	key := reqElements[1].(*respvalues.BasicBulkStr)

	expiryMSec := int64(0)

	if cmd == "EXPIRE" || cmd == "PEXPIRE" {
		expiryStr := reqElements[2].(*respvalues.BasicBulkStr).String()
		expiryParsed, expiryParseError := strconv.ParseInt(expiryStr, 10, 32)
		expiryMSec = expiryParsed

		if expiryParseError != nil {
			return respvalues.NewErr("EXPIRE/PEXPIRE/PERSIST: invalid expiry specified")
		}
		if expiryMSec < 0 {
			expiryMSec = -1
		}
		if cmd == "EXPIRE" {
			expiryMSec *= 1000
		}
	} else if cmd == "PERSIST" {
		expiryMSec = -1
	}

	keyHash := ringmgr.CalcKeyHash(key)
	bucketIdx := keyHash % uint32(len(buckets))
	bucketApi := buckets[bucketIdx]
	result := bucketApi.RunCmd(bucket.NewCmdExpire(key.String(), int64(expiryMSec)))

	return result
}

func (this *ExpireHandler) Register(registry map[string]CommandHandler) {
	registry["EXPIRE"] = this
	registry["PEXPIRE"] = this
	registry["PERSIST"] = this
}

func NewExpireHandler(ctx actor.Ctx, ringMgr ringmgr.RingMgr) CommandHandler {
	h := new(ExpireHandler)
	h.ctx = ctx
	h.ringMgr = ringMgr
	return h
}
