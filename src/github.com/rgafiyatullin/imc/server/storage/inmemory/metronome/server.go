package metronome

import (
	"container/list"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"time"
)

const TickDurationNanos = 10 * 1000 * 1000 // 10ms tick

type Tick interface {
	CurrentTickIdx() int64
}

type tick struct {
	t int64
}

func NewTick(idx int64) Tick {
	tick := new(tick)
	tick.t = idx
	return tick
}

func (this *tick) CurrentTickIdx() int64 {
	return this.t
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
	ctx         actor.Ctx
	initTime    time.Time
	inChans     *inChans
	joiners     *list.List
	subscribers map[chan<- Tick]bool
}

func metronomeEnterLoop(ctx actor.Ctx, inChans *inChans) {
	ctx.Log().Info("enter loop")
	state := new(state)
	state.ctx = ctx
	state.initTime = time.Now()
	state.inChans = inChans
	state.joiners = list.New()
	state.subscribers = make(map[chan<- Tick]bool)
	state.loop()
}

func (this *state) loop() {
	defer this.releaseJoiners()
	ticker := time.NewTicker(time.Nanosecond * TickDurationNanos)

	for {
		select {
		case join := <-this.inChans.join:
			this.joiners.PushBack(join)

		case subs := <-this.inChans.subsMgmt:
			this.handleSubsReq(subs)

		case tick := <-ticker.C:
			this.handleTick(tick)
		}
	}
}

func (this *state) releaseJoiners() {
	this.ctx.Log().Debug("releasing joiners...")
	join.ReleaseJoiners(this.joiners)
	this.joiners = list.New()
}

func (this *state) handleTick(t time.Time) {
	elapsed := t.Sub(this.initTime)
	tickIdx := int64(elapsed.Nanoseconds() / TickDurationNanos)
	//this.ctx.Log().Debug("tick #%d", tickIdx)

	tick := NewTick(tickIdx)
	for k, _ := range this.subscribers {
		k <- tick
	}
}

func (this *state) handleSubsReq(subs subsReq) {
	switch subs.(type) {
	case *subsReqSubscribe:
		this.ctx.Log().Debug("subscribe: %+v", subs.Chan())
		this.subscribers[subs.Chan()] = true
		subs.ReplyTo() <- true

	case *subsReqUnsubscribe:
		this.ctx.Log().Debug("unsubscribe: %+v", subs.Chan())
		delete(this.subscribers, subs.Chan())
		subs.ReplyTo() <- true

	default:
		this.ctx.Log().Debug("unexpected subs request: %v", subs)
	}
}
