package persistent

import "github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"

type nilPersister struct {
	restoreChan chan RestoreMsg
}

func CreateNilPersister() Persister {
	p := new(nilPersister)
	p.restoreChan = make(chan RestoreMsg, 1)
	p.restoreChan <- NewRestoreComplete()
	return p
}

func (this *nilPersister) Restore() <-chan RestoreMsg {
	return this.restoreChan
}
func (this *nilPersister) Save(key string, value data.Value) {}
