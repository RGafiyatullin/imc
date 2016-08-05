package netsrv

import (
	"github.com/rgafiyatullin/imc/protocol/resp/server"
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"github.com/rgafiyatullin/imc/server/actor"
	"net"
	"fmt"
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
			this.actorCtx.Log().Warning("error reading cmd: %v", err)
			break
		}
		this.processRequest(cmd)
	}
}

func (this *connectionState) processRequest(req *types.BasicArr) {
	elements := req.Elements()

	var resp types.BasicType = nil

	if len(elements) == 0 {
		resp = types.NewErr("Malformed command (0 parts)")
	} else {
		switch elements[0].(type) {
		case *types.BasicBulkStr:
			cmdName := elements[0].(*types.BasicBulkStr).String()
			switch cmdName {
			case "PING":
				resp = types.NewStr("PONG")
			default:
				resp = types.NewErr(fmt.Sprintf("Unknown command '%s'", cmdName))
			}

		default:
			resp = types.NewErr("Malformed command (expected first element to be a bulkStr)")
		}
	}

	this.actorCtx.Log().Debug("processRequest [req: %s; resp: %s]", req.ToString(), resp.ToString())
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
	this.actorCtx.Log().Warning("read error: %+v", err)
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
