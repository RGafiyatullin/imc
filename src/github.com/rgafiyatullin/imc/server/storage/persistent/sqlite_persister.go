package persistent

import (
	"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
)

const SaveChanBufSize = 64
const RestoreChanBufSize = 64

type channels struct {
	restore chan RestoreMsg
	save    chan SaveReq
}

type sqlitePersister struct {
	chans *channels
}

func (this *sqlitePersister) Restore() <-chan RestoreMsg {
	return this.chans.restore
}

func (this *sqlitePersister) Save(key string, value data.Value) {}

func StartSqlitePersister(actorCtx actor.Ctx, file string) Persister {
	p := new(sqlitePersister)
	chans := new(channels)
	chans.restore = make(chan RestoreMsg, RestoreChanBufSize)
	chans.save = make(chan SaveReq, SaveChanBufSize)
	p.chans = chans

	go persisterEnterLoop(actorCtx, chans, file)

	return p
}

const stRestoring = 0
const stSaving = 1

type state struct {
	actorCtx actor.Ctx
	chans    *channels
	file     string
	status   int
	db       *sql.DB
}

func (this *state) init(actorCtx actor.Ctx, chans *channels, file string) {
	this.actorCtx = actorCtx
	this.chans = chans
	this.file = file
	this.status = stRestoring
	this.actorCtx.Log().Info("init [file: '%s']", file)

	this.initSqlite()
}

func (this *state) initSqlite() {
	db, err := sql.Open("sqlite3", ":memory")
	if err != nil {
		this.actorCtx.Log().Fatal("Failed to open database [%s]: %v", this.file, err)
		this.actorCtx.Log().Flush()
		this.actorCtx.Halt(2, "SqlitePersister: Failed to open database")
	} else {
		this.db = db
	}
}

func (this *state) loop() {
	for {
		switch this.status {
		//case stExpectRestoreReq:
		//	this.whenExpectRestoreReq()
		case stRestoring:
			this.whenRestoring()
		case stSaving:
			this.whenSaving()
		}
	}
}

func persisterEnterLoop(actorCtx actor.Ctx, chans *channels, file string) {
	s := new(state)
	s.init(actorCtx, chans, file)
	s.loop()
}
