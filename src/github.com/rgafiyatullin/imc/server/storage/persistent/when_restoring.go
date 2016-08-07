package persistent

func (this *state) whenRestoring() {
	this.actorCtx.Log().Warning("whenRestoring: NOT IMPLEMENTED")
	this.status = stSaving
}
