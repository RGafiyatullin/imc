package ringmgr

import (
	"container/list"
	"fmt"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/metronome"
	"hash/crc32"
)

type RingMgr interface {
	Join() join.Awaitable
	QueryBuckets() []bucket.Bucket
}

type inChans struct {
	join         chan chan<- bool
	queryBuckets chan chan<- []bucket.Bucket
}
type ringmgr struct {
	chans *inChans
}

func (this *ringmgr) Join() join.Awaitable {
	ch := join.NewClientChan()
	this.chans.join <- ch
	return join.New(ch)
}

func (this *ringmgr) QueryBuckets() []bucket.Bucket {
	ch := make(chan []bucket.Bucket, 1)
	this.chans.queryBuckets <- ch
	return <-ch
}

func StartRingMgr(ctx actor.Ctx, config config.Config, metronome metronome.Metronome) RingMgr {
	chans := new(inChans)
	chans.join = join.NewServerChan()
	chans.queryBuckets = make(chan chan<- []bucket.Bucket, 32)
	m := new(ringmgr)
	m.chans = chans

	go ringMgrEnterLoop(ctx, config, chans, metronome)

	return m
}

func CalcKeyHash(key *respvalues.BasicBulkStr) uint32 {
	return crc32.ChecksumIEEE(key.Bytes())
}

type state struct {
	ctx     actor.Ctx
	chans   *inChans
	config  config.Config
	joiners *list.List
	buckets []bucket.Bucket
}

func (this *state) init(ctx actor.Ctx, config config.Config, chans *inChans, metronome metronome.Metronome) {
	this.joiners = list.New()
	this.ctx = ctx
	this.chans = chans
	this.config = config

	ringSize := this.config.Storage().RingSize()

	this.buckets = make([]bucket.Bucket, ringSize)
	var bucketIdx uint
	for bucketIdx = 0; bucketIdx < ringSize; bucketIdx++ {

		this.buckets[bucketIdx] = bucket.StartBucket(
			this.ctx.NewChild(fmt.Sprintf("bucket-%d", bucketIdx)),
			bucketIdx, config, metronome)
	}

	this.ctx.Log().Debug("init [ring-size: %d]", ringSize)

}

func (this *state) loop() {
	defer this.releaseJoiners()

	for {
		select {
		case join := <-this.chans.join:
			this.joiners.PushBack(join)
		case qb := <-this.chans.queryBuckets:
			qb <- this.buckets
		}
	}
}

func ringMgrEnterLoop(ctx actor.Ctx, config config.Config, chans *inChans, metronome metronome.Metronome) {
	ctx.Log().Info("enter loop")
	state := new(state)

	state.init(ctx, config, chans, metronome)
	state.loop()
}

func (this *state) releaseJoiners() {
	this.ctx.Log().Debug("releasing joiners...")
	join.ReleaseJoiners(this.joiners)
	this.joiners = list.New()
}
