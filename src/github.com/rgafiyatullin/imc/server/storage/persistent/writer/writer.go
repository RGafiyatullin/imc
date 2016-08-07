package writer

import (
	"database/sql"
	"github.com/rgafiyatullin/imc/server/actor"
)

type Writer interface{}

type chans struct{}
type writer struct {
	chans *chans
}

type state struct {
	ctx   actor.Ctx
	dsn   string
	chans *chans
	db    *sql.DB
}

func StartWriter(ctx actor.Ctx, dsn string) Writer {
	chans := &chans{}
	go enterLoop(ctx, dsn, chans)
	return &writer{
		chans: chans,
	}
}

func enterLoop(ctx actor.Ctx, dsn string, chans *chans) {
	state := &state{
		ctx:   ctx,
		dsn:   dsn,
		chans: chans,
	}
	state.init()
	state.loop()
}

func (this *state) init() {
	db, err := sql.Open("sqlite3", this.dsn)
	if err != nil {
		this.ctx.Log().Error("Failed to open database [%s]: %v", this.dsn, err)
		this.db = nil
	} else {
		this.db = db
	}
}

func (this *state) loop() {}
