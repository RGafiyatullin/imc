package netsrv

import (
	_ "github.com/rgafiyatullin/imc/protocol/resp/constatns"
	protoReader "github.com/rgafiyatullin/imc/protocol/resp/server"
	"github.com/rgafiyatullin/imc/protocol/resp/types/request"
	"github.com/rgafiyatullin/imc/server/actor"
	"net"
)

const ReadBufSize int = 10

type Connection interface {
	acceptorId() int
	connectionId() int
}

type connection struct {
	aId_ int
	cId_ int
}

type connectionState struct {
	actorCtx   actor.Ctx
	readBuf    []byte
	readerCtx  protoReader.Context
	sock       net.Conn
	cId        int
	aId        int
	closedChan chan<- *ClosedInfo
}

func (this *connectionState) init(ctx actor.Ctx, aId int, cId int, sock net.Conn, closedChan chan<- *ClosedInfo) {
	this.actorCtx = ctx
	this.actorCtx.Log().Debug("init")

	this.readBuf = make([]byte, ReadBufSize)
	this.sock = sock
	this.cId = cId
	this.aId = aId
	this.closedChan = closedChan
	this.readerCtx = protoReader.New()
}

func (this *connectionState) loop() {
	defer this.onClosed()

	for {
		bytesRead, err := this.sock.Read(this.readBuf)
		if err != nil {
			this.onReadError(err)
			return
		} else {
			chunk := make([]byte, bytesRead)
			copy(chunk, this.readBuf[:bytesRead])

			this.actorCtx.Log().Debug("adding chunk: %v", chunk)

			this.readerCtx.AddChunk(chunk)
			for this.readerCtx.HasRequest() {
				this.actorCtx.Log().Debug("has request")
				req := this.readerCtx.FetchRequest()
				this.processRequest(req)
			}
		}
	}
}

func (this *connectionState) processRequest(req request.Request) {
	this.actorCtx.Log().Debug("processRequest [req: %+v]", req)
}

func (this *connectionState) onClosed() {
	this.actorCtx.Log().Debug("closed")
	this.sock.Close()

	closedInfo := new(ClosedInfo)
	closedInfo.acceptorId = this.aId
	closedInfo.connectionId = this.cId
	this.closedChan <- closedInfo
}

func (this *connectionState) onReadError(err error) {
	this.actorCtx.Log().Warn("read error: %+v", err)
}

func (this *connection) acceptorId() int {
	return this.aId_
}
func (this *connection) connectionId() int {
	return this.cId_
}

func StartConnection(ctx actor.Ctx, aId int, cId int, sock net.Conn, closedChan chan<- *ClosedInfo) Connection {
	conn := new(connection)
	conn.aId_ = aId
	conn.cId_ = cId
	go connectionEnterLoop(ctx, aId, cId, sock, closedChan)
	return conn
}

func connectionEnterLoop(ctx actor.Ctx, aId int, cId int, sock net.Conn, closedChan chan<- *ClosedInfo) {
	ctx.Log().Debug("entering loop...", aId, cId, sock)

	state := new(connectionState)
	state.init(ctx, aId, cId, sock, closedChan)
	state.loop()
}
