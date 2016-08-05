package storage

import (
	"container/list"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/metronome"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
)

type Sup interface {
	Join() join.Awaitable
	QueryMetronome() metronome.Metronome
	QueryRingMgr() ringmgr.RingMgr
}

type sup struct {
	chans *inChans
}

func (this *sup) Join() join.Awaitable {
	// TODO: to be implemented
	return join.NewStub()
}

func (this *sup) QueryMetronome() metronome.Metronome {
	ch := make(chan metronome.Metronome, 1)
	this.chans.queryMetronome <- ch
	return <-ch
}

func (this *sup) QueryRingMgr() ringmgr.RingMgr {
	ch := make(chan ringmgr.RingMgr, 1)
	this.chans.queryRingmgr <- ch
	return <-ch
}

func StartSup(ctx actor.Ctx, config config.Config) Sup {
	chans := new(inChans)
	chans.queryMetronome = make(chan chan<- metronome.Metronome)
	chans.queryRingmgr = make(chan chan<- ringmgr.RingMgr)
	chans.join = join.NewServerChan()
	sup := new(sup)
	sup.chans = chans

	go supEnterLoop(ctx, config, chans)
	return sup
}

type inChans struct {
	join           chan chan<- bool
	queryMetronome chan chan<- metronome.Metronome
	queryRingmgr   chan chan<- ringmgr.RingMgr
}

type state struct {
	ctx     actor.Ctx
	config  config.Config
	chans   *inChans
	joiners *list.List
}

func (this *state) loop() {
	defer this.releaseJoiners()

	metronome := metronome.StartMetronome(this.ctx.NewChild("metronome"))
	ringmgr := ringmgr.StartRingMgr(this.ctx.NewChild("ring_mgr"), this.config, metronome)

	for {
		select {
		case join := <-this.chans.join:
			this.joiners.PushBack(join)
		case qm := <-this.chans.queryMetronome:
			qm <- metronome
		case qrm := <-this.chans.queryRingmgr:
			qrm <- ringmgr
		}
	}
}

func supEnterLoop(ctx actor.Ctx, config config.Config, chans *inChans) {
	ctx.Log().Info("enter loop")
	state := new(state)
	state.joiners = list.New()
	state.chans = chans
	state.ctx = ctx
	state.config = config
	state.loop()
}

func (this *state) releaseJoiners() {
	this.ctx.Log().Debug("releasing joiners...")
	join.ReleaseJoiners(this.joiners)
	this.joiners = list.New()
}
