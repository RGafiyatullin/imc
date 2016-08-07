package persistent

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
	"github.com/rgafiyatullin/imc/server/storage/persistent/writer"
)

const SaveChanBufSize = 64
const RestoreChanBufSize = 64

type channels struct {
	restore chan RestoreMsg
}

type sqlitePersister struct {
	chans *channels
}

func (this *sqlitePersister) Restore() <-chan RestoreMsg {
	return this.chans.restore
}

func (this *sqlitePersister) Save(key string, value data.Value) {}

func StartSqlitePersister(actorCtx actor.Ctx, file string) Persister {
	chans := &channels{
		restore: make(chan RestoreMsg, RestoreChanBufSize),
	}
	go persisterEnterLoop(actorCtx, chans, file)

	return &sqlitePersister{chans: chans}
}

const stRestoring = 0
const stSaving = 1

const ktString = 0
const ktList = 1
const ktMap = 2

type state struct {
	actorCtx actor.Ctx
	chans    *channels
	dsn      string
	status   int
	db       *sql.DB
	writer   writer.Writer
}

func (this *state) init(actorCtx actor.Ctx, chans *channels, file string) {
	this.actorCtx = actorCtx
	this.chans = chans
	this.dsn = "file:" + file
	this.status = stRestoring
	this.actorCtx.Log().Info("init [file: '%s']", file)

	this.initSqlite()
}

func (this *state) initSqlite() {
	db, err := sql.Open("sqlite3", this.dsn)
	if err != nil {
		this.actorCtx.Log().Fatal("Failed to open database [%s]: %v", this.dsn, err)
		this.actorCtx.Log().Flush()
		this.actorCtx.Halt(2, "SqlitePersister: Failed to open database")
	} else {
		this.db = db
		this.ensureSqliteSchema()
	}
}

func (this *state) ensureSqliteSchema() {
	_, keysErr := this.db.Exec("CREATE TABLE keys (" +
		"k varchar(1024) NOT NULL, " +
		"v blob, " +
		"t tinyint, " +
		"e integer, " +
		"PRIMARY KEY(k))")
	this.actorCtx.Log().Debug("keysErr: %v", keysErr)

	_, listIdErr := this.db.Exec("CREATE TABLE ids_seq (id integer NOT NULL, v tinyint, PRIMARY KEY(id))")

	this.actorCtx.Log().Debug("listIdErr: %v", listIdErr)

	_, listsErr := this.db.Exec("CREATE TABLE lists (" +
		"id integer NOT NULL, " +
		"list_id bigint NOT NULL, " +
		"idx int NOT NULL, " +
		"value blob NOT NULL," +
		"PRIMARY KEY(id))")
	this.actorCtx.Log().Debug("listsErr: %v", listsErr)

	_, mapsErr := this.db.Exec("CREATE TABLE maps (" +
		"id integer NOT NULL, " +
		"map_id bigint NOT NULL, " +
		"k varchar(1024), " +
		"value BLOB, " +
		"PRIMARY KEY(id))")
	this.actorCtx.Log().Debug("mapsErr: %v", mapsErr)
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
