package persistent

import "time"

func (this *state) startWriter() {
	this.actorCtx.Log().Error("startWriter: NOT IMPLEMENTED")
}

func (this *state) whenSaving() {
	this.actorCtx.Log().Warning("whenSaving: NOT IMPLEMENTED")
	time.Sleep(10 * time.Second)
}
