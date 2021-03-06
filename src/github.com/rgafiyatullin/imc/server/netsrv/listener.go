package netsrv

import (
	"container/list"
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"net"
	"time"
)

type listenerChannels struct {
	acceptedChan chan *AcceptedInfo
	closedChan   chan *ClosedInfo
	joinChan     chan chan<- bool
}

// Handle for Listener-actor
type Listener interface {
	Join() join.Awaitable
}

type listener struct {
	chans *listenerChannels
}

func (srv *listener) Join() join.Awaitable {
	ch := join.NewClientChan()
	srv.chans.joinChan <- ch
	return join.New(ch)
}

type AcceptedInfo struct {
	connection Connection
}
type ClosedInfo struct {
	acceptorId   int
	connectionId int
}

type srvState struct {
	actorCtx   actor.Ctx
	connsCount int
	ringMgr    ringmgr.RingMgr
	config     config.Config
	lSock      net.Listener
	chans      *listenerChannels
	joiners    *list.List
}

// Start a new Listener actor.
//
// Listener actor binds the interface according to the provided config and serves the Redis protocol on it.
//
// All the requests are passed to the provided RingMgr.
func StartListener(ctx actor.Ctx, config config.Config, ringMgr ringmgr.RingMgr) (Listener, error) {
	lSock, listenErr := net.Listen("tcp", config.Net().BindSpec())
	if listenErr != nil {
		ctx.Log().Fatal("Failed to bind '%s': %v", config.Net().BindSpec(), listenErr)
		ctx.Log().Flush()
		ctx.Halt(1, "netsrv.StartListener: bind error")
		return nil, listenErr
	} else {
		ctx.Log().Info("Bound '%s'", config.Net().BindSpec())
	}

	srv := new(listener)
	srv.chans = new(listenerChannels)
	srv.chans.acceptedChan = make(chan *AcceptedInfo)
	srv.chans.closedChan = make(chan *ClosedInfo)
	srv.chans.joinChan = join.NewServerChan()

	go listenerEnterLoop(ctx, lSock, srv.chans, ringMgr, config)

	return srv, nil
}

func listenerEnterLoop(
	actorCtx actor.Ctx, lSock net.Listener, chans *listenerChannels,
	ringMgr ringmgr.RingMgr, config config.Config) {

	state := new(srvState)
	state.init(actorCtx, lSock, chans, ringMgr, config)
	state.startAcceptors()
	state.listenerLoop()
}

func (this *srvState) init(
	actorCtx actor.Ctx, lSock net.Listener, chans *listenerChannels,
	ringMgr ringmgr.RingMgr, config config.Config) {

	this.actorCtx = actorCtx
	this.ringMgr = ringMgr
	this.config = config
	this.lSock = lSock
	this.chans = chans
	this.connsCount = 0
	this.joiners = list.New()
}

func (this *srvState) startAcceptors() {
	for idx := 0; idx < this.config.Net().AcceptorsCount(); idx++ {
		this.startAcceptor(idx)
	}
}

func (this *srvState) startAcceptor(idx int) {
	this.actorCtx.Log().Debug("Starting acceptor #%v", idx)
	childCtx := this.actorCtx.NewChild(fmt.Sprintf("acceptor-%v", idx))
	StartAcceptor(childCtx, idx, this.lSock, this.ringMgr, this.config, this.chans.acceptedChan, this.chans.closedChan)
}

func (this *srvState) listenerLoop() {
	defer this.releaseJoiners()
	metricsTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-this.chans.acceptedChan:
			this.connsCount++
			this.actorCtx.Metrics().ReportConnUp()

		case <-this.chans.closedChan:
			this.connsCount--
			this.actorCtx.Metrics().ReportConnDn()

		case <-metricsTicker.C:
			//this.actorCtx.Log().Debug("ConnsCount: %d", this.connsCount)
			this.actorCtx.Metrics().ReportConnCount(this.connsCount)

		case join := <-this.chans.joinChan:
			this.joiners.PushBack(join)
		}
	}
}

func (this *srvState) releaseJoiners() {
	this.actorCtx.Log().Debug("releasing joiners...")
	for element := this.joiners.Front(); element != nil; element = element.Next() {
		element.Value.(chan<- bool) <- true
	}
	this.joiners = list.New()
}
