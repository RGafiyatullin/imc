package netsrv

import (
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"net"
)

type Acceptor interface{}
type acceptor struct{}

type acceptorState struct {
	ctx          actor.Ctx
	ringMgr      ringmgr.RingMgr
	acceptorId   int
	lSock        net.Listener
	acceptedChan chan<- *AcceptedInfo
	closedChan   chan<- *ClosedInfo
}

func StartAcceptor(actorCtx actor.Ctx, acceptorId int, lSock net.Listener, ringMgr ringmgr.RingMgr, acceptedChan chan<- *AcceptedInfo, closedChan chan<- *ClosedInfo) Acceptor {
	go acceptorEnterLoop(actorCtx, acceptorId, lSock, ringMgr, acceptedChan, closedChan)
	acceptor := new(acceptor)
	return acceptor
}

func acceptorEnterLoop(ctx actor.Ctx, acceptorId int, lSock net.Listener, ringMgr ringmgr.RingMgr, acceptedChan chan<- *AcceptedInfo, closedChan chan<- *ClosedInfo) {
	ctx.Log().Debug("entering loop...")

	state := new(acceptorState)
	state.init(ctx, acceptorId, lSock, ringMgr, acceptedChan, closedChan)

	state.loop()
}

func (this *acceptorState) init(ctx actor.Ctx, acceptorId int, lSock net.Listener, ringMgr ringmgr.RingMgr, acceptedChan chan<- *AcceptedInfo, closedChan chan<- *ClosedInfo) {
	this.acceptorId = acceptorId
	this.ctx = ctx
	this.ringMgr = ringMgr
	this.acceptedChan = acceptedChan
	this.closedChan = closedChan
	this.lSock = lSock
}

func (this *acceptorState) loop() {
	idx := 0
	for {
		sock, err := this.lSock.Accept()

		if err != nil {
			this.ctx.Log().Error("error on accept [%v]:  %v", idx, err)
		} else {
			this.ctx.Log().Debug("accepted: [%v] %v", idx, sock)
			connection := StartConnection(
				this.ctx.NewChild(fmt.Sprintf("conn-%v", idx)),
				this.acceptorId, idx, sock, this.ringMgr, this.closedChan)
			acceptedInfo := new(AcceptedInfo)
			acceptedInfo.connection = connection
			this.acceptedChan <- acceptedInfo
		}

		idx++
	}
}
