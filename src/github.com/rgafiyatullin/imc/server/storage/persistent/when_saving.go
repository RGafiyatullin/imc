package persistent

import (
	"github.com/rgafiyatullin/imc/server/storage/persistent/writer"
	"time"
)

func (this *state) startWriter() {
	this.writer = writer.StartWriter(this.actorCtx.NewChild("writer"), this.dsn)
}

func (this *state) whenSaving() {
	this.actorCtx.Log().Warning("whenSaving: NOT IMPLEMENTED")
	time.Sleep(10 * time.Second)
}
