package metronome

import "github.com/rgafiyatullin/imc/server/actor/join"

const MetronomeTickChanBufSize = 16

type Metronome interface {
	Subscribe(tickChan chan<- Tick)
	Unsubscribe(tickChan chan<- Tick)
	Join() join.Awaitable
}

func NewChan() chan Tick {
	return make(chan Tick, MetronomeTickChanBufSize)
}

func (this *metronome) Subscribe(tickChan chan<- Tick) {
	req := new(subsReqSubscribe)
	req.tickChan = tickChan
	replyChan := make(chan bool, 1)
	req.replyChan = replyChan

	this.inChans.subsMgmt <- req
	<-replyChan
}

func (this *metronome) Unsubscribe(tickChan chan<- Tick) {
	req := new(subsReqUnsubscribe)
	req.tickChan = tickChan
	replyChan := make(chan bool, 1)
	req.replyChan = replyChan

	this.inChans.subsMgmt <- req
	<-replyChan
}

func (this *metronome) Join() join.Awaitable {
	ch := join.NewClientChan()
	this.inChans.join <- ch
	return join.New(ch)
}

type subsReq interface {
	Chan() chan<- Tick
	ReplyTo() chan<- bool
}

type subsReqSubscribe struct {
	tickChan  chan<- Tick
	replyChan chan bool
}

func (this *subsReqSubscribe) Chan() chan<- Tick    { return this.tickChan }
func (this *subsReqSubscribe) ReplyTo() chan<- bool { return this.replyChan }

type subsReqUnsubscribe struct {
	tickChan  chan<- Tick
	replyChan chan bool
}

func (this *subsReqUnsubscribe) Chan() chan<- Tick    { return this.tickChan }
func (this *subsReqUnsubscribe) ReplyTo() chan<- bool { return this.replyChan }
