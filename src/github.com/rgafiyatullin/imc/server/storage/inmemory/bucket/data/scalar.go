package data

import "github.com/rgafiyatullin/imc/protocol/resp/respvalues"

type ScalarValue struct {
	value []byte
}

func NewScalar(value []byte) *ScalarValue {
	v := new(ScalarValue)
	v.value = make([]byte, len(value))
	copy(v.value, value)
	return v
}

func (this *ScalarValue) ToRESP() respvalues.BasicType {
	return respvalues.NewBulkStr(this.value)
}
