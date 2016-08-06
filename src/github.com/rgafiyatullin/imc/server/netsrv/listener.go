package netsrv

import (
	"container/list"
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"net"
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
	actorCtx       actor.Ctx
	ringMgr        ringmgr.RingMgr
	acceptorsCount int
	lSock          net.Listener
	chans          *listenerChannels
	joiners        *list.List
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

	go listenerEnterLoop(ctx, lSock, srv.chans, ringMgr)

	return srv, nil
}

func listenerEnterLoop(actorCtx actor.Ctx, lSock net.Listener, chans *listenerChannels, ringMgr ringmgr.RingMgr) {
	state := new(srvState)
	state.init(actorCtx, 10, lSock, chans, ringMgr)
	state.startAcceptors()
	state.listenerLoop()
}

func (this *srvState) init(actorCtx actor.Ctx, acceptorsCount int, lSock net.Listener, chans *listenerChannels, ringMgr ringmgr.RingMgr) {
	this.actorCtx = actorCtx
	this.ringMgr = ringMgr
	this.acceptorsCount = acceptorsCount
	this.lSock = lSock
	this.chans = chans
	this.joiners = list.New()
}

func (this *srvState) startAcceptors() {
	for idx := 0; idx < this.acceptorsCount; idx++ {
		this.startAcceptor(idx)
	}
}

func (this *srvState) startAcceptor(idx int) {
	this.actorCtx.Log().Debug("Starting acceptor #%v", idx)
	childCtx := this.actorCtx.NewChild(fmt.Sprintf("acceptor-%v", idx))
	StartAcceptor(childCtx, idx, this.lSock, this.ringMgr, this.chans.acceptedChan, this.chans.closedChan)
}

func (this *srvState) listenerLoop() {
	defer this.releaseJoiners()
	for {
		select {
		case accepted := <-this.chans.acceptedChan:
			this.actorCtx.Log().Debug("accepted %+v", *accepted)
		case closed := <-this.chans.closedChan:
			this.actorCtx.Log().Debug("closed %+v", *closed)
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
