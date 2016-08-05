package data

import (
	"github.com/rgafiyatullin/imc/protocol/resp/respvalues"
)

type Value interface {
	ToRESP() respvalues.BasicType
}
