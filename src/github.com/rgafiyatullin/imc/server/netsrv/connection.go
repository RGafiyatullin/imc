package netsrv

import (
	"fmt"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/protocol/resp/server"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/netsrv/commands"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/ringmgr"
	"net"
	"time"
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
	ringMgr    ringmgr.RingMgr
	cId        int
	aId        int
	closedChan chan<- *ClosedInfo
	handlers   map[string]commands.CommandHandler
}

func (this *connectionState) init(ctx actor.Ctx, aId int, cId int, sock net.Conn, ringMgr ringmgr.RingMgr, closedChan chan<- *ClosedInfo) {
	this.actorCtx = ctx
	this.actorCtx.Log().Debug("init")

	this.readBuf = make([]byte, ReadBufSize)
	this.sock = sock
	this.ringMgr = ringMgr
	this.cId = cId
	this.aId = aId
	this.closedChan = closedChan
	this.protocol = server.New(sock)

	this.initCommands()
}

func (this *connectionState) initCommands() {
	this.handlers = make(map[string]commands.CommandHandler)
	commands.NewPingHandler(this.actorCtx.NewChild("#PING")).Register(this.handlers)
	commands.NewGetHandler(this.actorCtx.NewChild("#GET"), this.ringMgr).Register(this.handlers)
	commands.NewSetHandler(this.actorCtx.NewChild("#SET"), this.ringMgr).Register(this.handlers)
	commands.NewDelHandler(this.actorCtx.NewChild("#DEL"), this.ringMgr).Register(this.handlers)
	commands.NewLPshFHandler(this.actorCtx.NewChild("#LPSHF"), this.ringMgr).Register(this.handlers)
	commands.NewLPshBHandler(this.actorCtx.NewChild("#LPSHB"), this.ringMgr).Register(this.handlers)
	commands.NewLPopFHandler(this.actorCtx.NewChild("#LPOPF"), this.ringMgr).Register(this.handlers)
	commands.NewLPopBHandler(this.actorCtx.NewChild("#LPOPB"), this.ringMgr).Register(this.handlers)
	commands.NewLNthHandler(this.actorCtx.NewChild("#LNTH"), this.ringMgr).Register(this.handlers)
}

func (this *connectionState) loop() {
	defer this.onClosed()

	for {
		cmd, err := this.protocol.NextCommand()

		cmdExecStart := time.Now()
		if err != nil {
			this.actorCtx.Log().Warning("error reading cmd: %v", err)
			break
		}
		this.processRequest(cmd)
		cmdExecElapsed := time.Since(cmdExecStart)
		this.actorCtx.Metrics().ReportCommandDuration(cmdExecElapsed)
	}
}

func (this *connectionState) processRequest(req *respvalues.BasicArr) {
	elements := req.Elements()

	var resp respvalues.BasicType = nil

	if len(elements) == 0 {
		resp = respvalues.NewErr("Malformed command (0 parts)")
	} else {
		switch elements[0].(type) {
		case *respvalues.BasicBulkStr:
			cmdName := elements[0].(*respvalues.BasicBulkStr).String()

			handler, ok := this.handlers[cmdName]
			if ok {
				resp = handler.Handle(req)
			} else {
				resp = respvalues.NewErr(fmt.Sprintf("Unknown command '%s'", cmdName))
			}

		default:
			resp = respvalues.NewErr("Malformed command (expected first element to be a bulkStr)")
		}
	}

	//this.actorCtx.Log().Debug("processRequest [req: %s; resp: %s]", req.ToString(), resp.ToString())

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

func StartConnection(ctx actor.Ctx, aId int, cId int, sock net.Conn, ringMgr ringmgr.RingMgr, closedChan chan<- *ClosedInfo) Connection {
	conn := new(connection)
	conn.aId_ = aId
	conn.cId_ = cId

	go connectionEnterLoop(ctx, aId, cId, sock, ringMgr, closedChan)
	return conn
}

func connectionEnterLoop(ctx actor.Ctx, aId int, cId int, sock net.Conn, ringMgr ringmgr.RingMgr, closedChan chan<- *ClosedInfo) {
	ctx.Log().Debug("entering loop...", aId, cId, sock)

	state := new(connectionState)
	state.init(ctx, aId, cId, sock, ringMgr, closedChan)
	state.loop()
}
