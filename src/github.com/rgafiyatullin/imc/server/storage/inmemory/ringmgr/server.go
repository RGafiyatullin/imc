package ringmgr

import (
	"container/list"
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket"
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
	ch := make(chan []bucket.Bucket)
	this.chans.queryBuckets <- ch
	return <-ch
}

func StartRingMgr(ctx actor.Ctx, config config.Config) RingMgr {
	chans := new(inChans)
	chans.join = join.NewServerChan()
	m := new(ringmgr)
	m.chans = chans

	go ringMgrEnterLoop(ctx, config, chans)

	return m
}

type state struct {
	ctx     actor.Ctx
	chans   *inChans
	config  config.Config
	joiners *list.List
	buckets []bucket.Bucket
}

func (this *state) init(ctx actor.Ctx, config config.Config, chans *inChans) {
	this.joiners = list.New()
	this.ctx = ctx
	this.chans = chans
	this.config = config

	ringSize := this.config.Storage().RingSize()

	this.buckets = make([]bucket.Bucket, ringSize)
	var bucketIdx uint
	for bucketIdx = 0; bucketIdx < ringSize; bucketIdx++ {

		this.buckets[bucketIdx] = bucket.StartBucket(
			this.ctx.NewChild(fmt.Sprintf("bucket-%d", bucketIdx)), bucketIdx, config)
	}

	this.ctx.Log().Debug("init [ring-size: %d]", ringSize)

}

func (this *state) loop() {
	defer this.releaseJoiners()

	for {
		select {
		case join := <-this.chans.join:
			this.joiners.PushBack(join)
		}
	}
}

func ringMgrEnterLoop(ctx actor.Ctx, config config.Config, chans *inChans) {
	ctx.Log().Info("enter loop")
	state := new(state)

	state.init(ctx, config, chans)
	state.loop()
}

func (this *state) releaseJoiners() {
	this.ctx.Log().Debug("releasing joiners...")
	join.ReleaseJoiners(this.joiners)
	this.joiners = list.New()
}
