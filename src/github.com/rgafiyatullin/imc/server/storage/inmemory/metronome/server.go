package metronome

import (
	"container/list"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
)

type Tick interface {
	t() uint64
}

type tick struct {
	t_ uint64
}

func (this *tick) t() uint64 {
	return this.t_
}

type inChans struct {
	subsMgmt chan subsReq
	join     chan chan<- bool
}

type metronome struct {
	inChans *inChans
}

func StartMetronome(ctx actor.Ctx) Metronome {
	m := new(metronome)
	inChans := new(inChans)
	inChans.subsMgmt = make(chan subsReq)
	inChans.join = join.NewServerChan()
	m.inChans = inChans

	go metronomeEnterLoop(ctx, inChans)

	return m
}

type state struct {
	ctx     actor.Ctx
	inChans *inChans
	joiners *list.List
}

func metronomeEnterLoop(ctx actor.Ctx, inChans *inChans) {
	ctx.Log().Info("enter loop")
	state := new(state)
	state.ctx = ctx
	state.inChans = inChans
	state.joiners = list.New()
	state.loop()
}

func (this *state) loop() {
	defer this.releaseJoiners()
	for {
		select {
		case join := <-this.inChans.join:
			this.joiners.PushBack(join)

		case subs := <-this.inChans.subsMgmt:
			this.handleSubsReq(subs)
		}
	}
}

func (this *state) releaseJoiners() {
	this.ctx.Log().Debug("releasing joiners...")
	join.ReleaseJoiners(this.joiners)
	this.joiners = list.New()
}

func (this *state) handleSubsReq(subs subsReq) {
	switch subs.(type) {
	case *subsReqSubscribe:
		this.ctx.Log().Debug("subscribe: %+v", subs.Chan())
		subs.ReplyTo() <- true

	case *subsReqUnsubscribe:
		this.ctx.Log().Debug("unsubscribe: %+v", subs.Chan())
		subs.ReplyTo() <- true

	default:
		this.ctx.Log().Debug("unexpected subs request: %v", subs)
	}
}
