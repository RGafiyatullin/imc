package bucket

import (
	"container/list"
	"fmt"

	"errors"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/metronome"
	"github.com/rgafiyatullin/imc/server/storage/persistent"
	"path"
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
	ctx    actor.Ctx
	idx    uint
	config config.Config

	chans   *inChans
	joiners *list.List

	storage *storage

	metronome metronome.Metronome
	tickChan  <-chan metronome.Tick

	persister         persistent.Persister
	restoreInProgress bool
	restoreChan       <-chan persistent.RestoreMsg
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

	if config.Storage().PersistenceEnabled() {
		totalBucketsCount := config.Storage().RingSize()
		persisterFile := path.Join(
			config.Storage().PersistenceDirectory(),
			fmt.Sprintf("bucket-%d-%d.sqlite", totalBucketsCount, idx))
		this.persister = persistent.StartSqlitePersister(this.ctx.NewChild("sqlite_persister"), persisterFile)

	} else {
		this.persister = persistent.CreateNilPersister()
	}
	this.restoreInProgress = true
	this.restoreChan = this.persister.Restore()

	this.ctx.Log().Debug("init")
}

func (this *state) loop() {
	for {
		select {
		case tick := <-this.tickChan:
			this.storage.tickIdx = tick.CurrentTickIdx()
			this.storage.PurgeTimedOut()
			this.storage.ReportStats()

		case join := <-this.chans.join:
			this.joiners.PushBack(join)

		case restoreMsg := <-this.restoreChan:
			this.handleRestoreMsg(restoreMsg)

		case cmdReq := <-this.chans.cmd:
			result, err := this.handleCommandRequest(cmdReq)
			if err != nil {
				response := respvalues.NewErr(fmt.Sprintf("%v", err))
				cmdReq.ReplyTo() <- response
			} else {
				cmdReq.ReplyTo() <- result
			}

		}
	}
}

func (this *state) handleRestoreMsg(restoreMsg persistent.RestoreMsg) {
	if this.restoreInProgress {
		switch restoreMsg.(type) {
		case (*persistent.RestoreComplete):
			this.ctx.Log().Info("restore complete: %+v", restoreMsg)
			this.restoreInProgress = false
		default:
			this.ctx.Log().Info("restore: received %+v", restoreMsg)
		}

	} else {
		this.ctx.Log().Warning("received unexpected RestoreMsg(%+v)", restoreMsg)
	}
}

func (this *state) handleCommandRequest(cmdReq CmdReq) (respvalues.RESPValue, error) {
	if this.restoreInProgress {
		return nil, errors.New("RESTORE_IN_PROGRESS")
	} else {
		return this.storage.handleCommand(cmdReq.Cmd())
	}
}

func bucketEnterLoop(ctx actor.Ctx, idx uint, config config.Config, chans *inChans, metronome metronome.Metronome) {
	state := new(state)

	state.init(ctx, idx, config, chans, metronome)
	state.loop()
}
