package persistent

import (
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
)

const SaveChanBufSize = 64

type inChannels struct {
	restore chan RestoreReq
	save    chan SaveReq
}

type sqlitePersister struct {
	chans *inChannels
}

func (this *sqlitePersister) Restore() <-chan RestoreMsg {
	return nil
}

func (this *sqlitePersister) Save(key string, value data.Value) {}

func StartSqlitePersister(actorCtx actor.Ctx, file string) Persister {
	p := new(sqlitePersister)
	chans := new(inChannels)
	chans.restore = make(chan RestoreReq, 1)
	chans.save = make(chan SaveReq, SaveChanBufSize)
	p.chans = chans

	go persisterEnterLoop(actorCtx, chans, file)

	return p
}

const stExpectRestoreReq = 0
const stRestoring = 1
const stSaving = 2

type state struct {
	actorCtx actor.Ctx
	chans    *inChannels
	file     string
	status   int
}

func (this *state) init(actorCtx actor.Ctx, chans *inChannels, file string) {
	this.actorCtx = actorCtx
	this.chans = chans
	this.file = file
	this.status = stExpectRestoreReq
}

func (this *state) loop() {
	for {
		switch this.status {
		case stExpectRestoreReq:
			this.whenExpectRestoreReq()
		case stRestoring:
			this.whenRestoring()
		case stSaving:
			this.whenSaving()
		}
	}
}

func persisterEnterLoop(actorCtx actor.Ctx, chans *inChannels, file string) {
	s := new(state)
	s.init(actorCtx, chans, file)
	s.loop()
}
