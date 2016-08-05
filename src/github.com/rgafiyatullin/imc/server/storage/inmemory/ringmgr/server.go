package ringmgr

import (
	"container/list"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
)

type RingMgr interface {
	Join() join.Awaitable
}

type inChans struct {
	join chan chan<- bool
}
type ringmgr struct {
	chans *inChans
}

func (this *ringmgr) Join() join.Awaitable {
	// TODO: to be implemented
	return join.NewStub()
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
	state.ctx = ctx
	state.chans = chans
	state.config = config
	state.joiners = list.New()
	state.loop()
}

func (this *state) releaseJoiners() {
	this.ctx.Log().Debug("releasing joiners...")
	join.ReleaseJoiners(this.joiners)
	this.joiners = list.New()
}
