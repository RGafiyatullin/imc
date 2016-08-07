package persistent

import "time"

func (this *state) whenRestoring() {
	this.actorCtx.Log().Warning("whenRestoring: NOT IMPLEMENTED")
	this.status = stSaving
	time.Sleep(5 * time.Second)
	this.chans.restore <- NewRestoreComplete()
}
