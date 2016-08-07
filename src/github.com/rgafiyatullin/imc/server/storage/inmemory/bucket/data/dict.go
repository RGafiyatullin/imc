package data

import (
	"container/list"
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
)

type DictValue struct {
	elements map[string]*ScalarValue
}

func (this *DictValue) ToRESP() respvalues.RESPValue {
	tuples := list.New()
	for k, v := range this.elements {
		tuple := list.New()
		tuple.PushBack(respvalues.NewStr(k))
		tuple.PushBack(v.ToRESP())
		tuples.PushBack(respvalues.NewArray(tuple))
	}
	return respvalues.NewArray(tuples)
}

func NewDict() *DictValue {
	dict := new(DictValue)
	dict.elements = make(map[string]*ScalarValue)
	return dict
}

func (this *DictValue) Set(key string, value []byte) (created bool) {
	_, existed := this.elements[key]
	this.elements[key] = NewScalar(value)
	created = !existed
	return created
}

func (this *DictValue) Get(key string) (value []byte, keyfound bool) {
	sc, found := this.elements[key]
	if !found {
		return nil, false
	} else {
		return sc.value, true
	}
}

func (this *DictValue) Del(key string) (existed bool, empty bool) {
	_, existed = this.elements[key]
	if existed {
		delete(this.elements, key)
	}
	empty = len(this.elements) == 0
	return existed, empty
}

func (this *DictValue) Keys() []string {
	size := len(this.elements)
	keys := make([]string, size)
	idx := 0
	for key, _ := range this.elements {
		keys[idx] = key
		idx++
	}
	return keys
}

func (this *DictValue) Values() [][]byte {
	size := len(this.elements)
	values := make([][]byte, size)
	idx := 0
	for _, sc := range this.elements {
		values[idx] = sc.value
		idx++
	}
	return values
}
