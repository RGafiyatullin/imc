package netsrv

import (
	"github.com/rgafiyatullin/imc/protocol/resp/server"
	"github.com/rgafiyatullin/imc/server/actor"
	"net"
	"github.com/rgafiyatullin/imc/protocol/resp/types"
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
	protocol   server.Context
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
	this.protocol = server.New(sock)
}

func (this *connectionState) loop() {
	defer this.onClosed()

	for {
		cmd, err := this.protocol.NextCommand()
		if err != nil {
			this.actorCtx.Log().Warn("error reading cmd: %v", err)
			break
		}
		//this.actorCtx.Log().Debug("command: %+v", cmd.ToString())
		this.processRequest(cmd)
	}
}

func (this *connectionState) processRequest(req *types.BasicArr) {
	this.actorCtx.Log().Debug("processRequest [req: %+v]", req.ToString())
	resp := types.NewErr("Not implemented. Coming soon :)")
	this.actorCtx.Log().Debug("About to write response...")
	this.protocol.Write(resp)
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
