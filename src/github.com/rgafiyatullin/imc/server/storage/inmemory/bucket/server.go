package bucket

import (
	"container/list"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/protocol/resp/types"
	"fmt"
)

type Bucket interface {
	Join() join.Awaitable
	RunCmd(cmd Cmd) types.BasicType
}

type bucket struct {
	chans *inChans
}

func (this *bucket) RunCmd(cmd Cmd) types.BasicType {
	ch := make(chan types.BasicType, 1)
	req := new(cmdReq)
	req.cmd = cmd
	req.replyTo = ch
	return <- ch
}

func (this *bucket) Join() join.Awaitable {
	ch := join.NewClientChan()
	this.chans.join <- ch
	return join.New(ch)
}

type CmdReq interface {
	ReplyTo() chan <- types.BasicType
	Cmd() Cmd
}
type cmdReq struct {
	replyTo chan types.BasicType
	cmd Cmd
}

func (this *cmdReq) ReplyTo() chan <- types.BasicType { return this.replyTo }
func (this *cmdReq) Cmd() Cmd { return this.cmd }

type inChans struct {
	join chan chan<- bool
	cmd chan CmdReq
}

func StartBucket(ctx actor.Ctx, idx uint, config config.Config) Bucket {
	chans := new(inChans)
	chans.join = join.NewServerChan()
	chans.cmd = make(chan CmdReq)
	bucket := new(bucket)
	bucket.chans = chans

	go bucketEnterLoop(ctx, idx, config, chans)

	return bucket
}

type state struct {
	ctx     actor.Ctx
	idx     uint
	config  config.Config
	chans   *inChans
	joiners *list.List
	storage *storage
}

func (this *state) init(ctx actor.Ctx, idx uint, config config.Config, chans *inChans) {
	this.ctx = ctx
	this.idx = idx
	this.config = config
	this.chans = chans
	this.joiners = list.New()
	this.storage = NewStorage()

	this.ctx.Log().Debug("init")
}

func (this *state) loop() {
	for {
		select {
		case join := <-this.chans.join:
			this.joiners.PushBack(join)
		case cmdReq := <- this.chans.cmd:
			result, err := this.storage.handleCommand(cmdReq.Cmd())
			if err != nil {
				response := types.NewErr(fmt.Sprintf("%v", err))
				cmdReq.ReplyTo() <- response
			} else {
				cmdReq.ReplyTo() <- result
			}

		}
	}
}

func bucketEnterLoop(ctx actor.Ctx, idx uint, config config.Config, chans *inChans) {
	state := new(state)

	state.init(ctx, idx, config, chans)
	state.loop()
}
