package persistent

import "time"

func (this *state) whenSaving() {
	this.actorCtx.Log().Warning("whenSaving: NOT IMPLEMENTED")
	time.Sleep(10 * time.Second)
}
