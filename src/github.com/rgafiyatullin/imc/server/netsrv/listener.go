package netsrv

import (
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"net"
)

type listenerChannels struct {
	acceptedChan   chan *AcceptedInfo
	closedChan     chan *ClosedInfo
	terminatedChan chan *TerminatedInfo
}

type Listener interface {
	Join()
}

type listener struct {
	chans *listenerChannels
}

func (srv *listener) Join() {
	<-srv.chans.terminatedChan
}

type AcceptedInfo struct {
	connection Connection
}
type ClosedInfo struct {
	acceptorId   int
	connectionId int
}
type TerminatedInfo struct{}

type srvState struct {
	actorCtx       actor.Ctx
	acceptorsCount int
	lSock          net.Listener
	chans          *listenerChannels
}

func StartListener(ctx actor.Ctx, addrSpec string) (Listener, error) {
	ctx.Log().Debug("Hello there!")
	lSock, listenErr := net.Listen("tcp", addrSpec)
	if listenErr != nil {
		return nil, listenErr
	}

	srv := new(listener)
	srv.chans = new(listenerChannels)
	srv.chans.acceptedChan = make(chan *AcceptedInfo)
	srv.chans.closedChan = make(chan *ClosedInfo)
	srv.chans.terminatedChan = make(chan *TerminatedInfo)

	go listenerEnterLoop(ctx, lSock, srv.chans)

	return srv, nil
}

func listenerEnterLoop(actorCtx actor.Ctx, lSock net.Listener, chans *listenerChannels) {
	state := new(srvState)
	state.init(actorCtx, 10, lSock, chans)
	state.startAcceptors()
	state.listenerLoop()
}

func (this *srvState) init(actorCtx actor.Ctx, acceptorsCount int, lSock net.Listener, chans *listenerChannels) {
	this.actorCtx = actorCtx
	this.acceptorsCount = acceptorsCount
	this.lSock = lSock
	this.chans = chans
}

func (this *srvState) startAcceptors() {
	for idx := 0; idx < this.acceptorsCount; idx++ {
		this.startAcceptor(idx)
	}
}

func (this *srvState) startAcceptor(idx int) {
	this.actorCtx.Log().Debug("Starting acceptor #%v", idx)
	childCtx := this.actorCtx.NewChild(fmt.Sprintf("acceptor-%v", idx))
	StartAcceptor(childCtx, idx, this.lSock, this.chans.acceptedChan, this.chans.closedChan)
}

func (this *srvState) listenerLoop() {
	for {
		select {
		case accepted := <-this.chans.acceptedChan:
			this.actorCtx.Log().Debug("accepted %+v", *accepted)
		case closed := <-this.chans.closedChan:
			this.actorCtx.Log().Debug("closed %+v", *closed)
		}
	}
}
