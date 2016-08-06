package bucket

import (
	"container/list"
	"fmt"

	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/metronome"
)

type Bucket interface {
	Join() join.Awaitable
	RunCmd(cmd Cmd) respvalues.RESPValue
}

type bucket struct {
	chans *inChans
}

func (this *bucket) RunCmd(cmd Cmd) respvalues.RESPValue {
	ch := make(chan respvalues.RESPValue, 1)
	req := new(cmdReq)
	req.cmd = cmd
	req.replyTo = ch
	this.chans.cmd <- req
	return <-ch
}

func (this *bucket) Join() join.Awaitable {
	ch := join.NewClientChan()
	this.chans.join <- ch
	return join.New(ch)
}

type CmdReq interface {
	ReplyTo() chan<- respvalues.RESPValue
	Cmd() Cmd
}
type cmdReq struct {
	replyTo chan respvalues.RESPValue
	cmd     Cmd
}

func (this *cmdReq) ReplyTo() chan<- respvalues.RESPValue { return this.replyTo }
func (this *cmdReq) Cmd() Cmd                             { return this.cmd }

type inChans struct {
	join chan chan<- bool
	cmd  chan CmdReq
}

func StartBucket(ctx actor.Ctx, idx uint, config config.Config, metronome metronome.Metronome) Bucket {
	chans := new(inChans)
	chans.join = join.NewServerChan()
	chans.cmd = make(chan CmdReq)
	bucket := new(bucket)
	bucket.chans = chans

	go bucketEnterLoop(ctx, idx, config, chans, metronome)

	return bucket
}

type state struct {
	ctx       actor.Ctx
	idx       uint
	config    config.Config
	chans     *inChans
	joiners   *list.List
	storage   *storage
	metronome metronome.Metronome
	tickChan  <-chan metronome.Tick
}

func (this *state) init(ctx actor.Ctx, idx uint, config config.Config, chans *inChans, m metronome.Metronome) {
	this.ctx = ctx
	this.idx = idx
	this.config = config
	this.chans = chans
	this.joiners = list.New()
	this.storage = NewStorage(ctx.NewChild("#storage"))
	this.metronome = m

	tickChan := metronome.NewChan()
	this.tickChan = tickChan
	this.metronome.Subscribe(tickChan)

	this.ctx.Log().Debug("init")
}

func (this *state) loop() {
	for {
		select {
		case tick := <-this.tickChan:
			this.storage.tickIdx = tick.CurrentTickIdx()
			this.storage.PurgeTimedOut()

		case join := <-this.chans.join:
			this.joiners.PushBack(join)
		case cmdReq := <-this.chans.cmd:
			result, err := this.storage.handleCommand(cmdReq.Cmd())
			if err != nil {
				response := respvalues.NewErr(fmt.Sprintf("%v", err))
				cmdReq.ReplyTo() <- response
			} else {
				cmdReq.ReplyTo() <- result
			}

		}
	}
}

func bucketEnterLoop(ctx actor.Ctx, idx uint, config config.Config, chans *inChans, metronome metronome.Metronome) {
	state := new(state)

	state.init(ctx, idx, config, chans, metronome)
	state.loop()
}
