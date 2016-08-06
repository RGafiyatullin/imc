package metronome

import "github.com/rgafiyatullin/imc/server/actor/join"

const MetronomeTickChanBufSize = 16

// Metronome actor handle
type Metronome interface {
	// Subscribe the channel for the ticks by this metronome
	//
	// NB. After the subscription the channel should be timely read in order not to block the whole metronome.
	// Hence the constant MetronomeTickChanBufSize.
	Subscribe(tickChan chan<- Tick)
	// Unsubscribe the channel from the ticks by this metronome
	//
	// NB. Asynchronous. I.e. we cannot guarantee synchronicity here: even though
	// the metronome actor won't send any ticks into the channel
	// after this call is dispatched; there cannot be the guarantee that
	// a tick won't be emitted between the last receive on the channel and this call.
	Unsubscribe(tickChan chan<- Tick)

	// See Joinable
	Join() join.Awaitable
}

// creates a new channel (used for tick-subscription; see Metronom.Subscribe and Metronom.Unsubscribe methods)
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
