package persistent

import "github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"

type Persister interface {
	Restore() <-chan RestoreMsg
	Save(key string, value data.Value)
}
