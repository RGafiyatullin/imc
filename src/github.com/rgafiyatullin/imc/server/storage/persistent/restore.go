package persistent

import "github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"

type RestoreMsg interface {
	IsComplete() bool
}

func NewRestoreComplete() *RestoreComplete {
	return &RestoreComplete{}
}

func NewRestoreString(k string, v *data.ScalarValue) *RestoreStringValue {
	return &RestoreStringValue{ Key: k, Value: v }
}

func NewRestoreList(k string, v *data.ListValue) *RestoreListValue {
	return &RestoreListValue{ Key: k, Value: v }
}

func NewRestoreDict(k string, v *data.DictValue) *RestoreDictValue {
	return &RestoreDictValue{ Key: k, Value: v }
}


type RestoreStringValue struct {
	Key   string
	Value *data.ScalarValue
}

type RestoreListValue struct {
	Key string
	Value *data.ListValue
}

type RestoreDictValue struct {
	Key string
	Value *data.DictValue
}

func (this *RestoreStringValue) IsComplete() bool { return false }
func (this *RestoreListValue) IsComplete() bool { return false }
func (this *RestoreDictValue) IsComplete() bool { return false }

type RestoreComplete struct{}

func (this *RestoreComplete) IsComplete() bool { return true }

