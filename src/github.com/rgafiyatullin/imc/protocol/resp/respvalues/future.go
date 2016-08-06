package respvalues

import (
	"fmt"
	"net/textproto"
)

func NewFuture() *RESPFuture {
	fut := new(RESPFuture)
	fut.resolved = nil
	fut.ch = make(chan RESPValue, 1)
	return fut
}

type RESPFuture struct {
	ch       chan RESPValue
	resolved RESPValue
}

func (this *RESPFuture) ToString() string {
	if this.resolved != nil {
		return fmt.Sprintf("F(%s)", this.resolved.ToString())
	} else {
		return "F(???)"
	}
}

func (this *RESPFuture) Write(to *textproto.Conn) {
	if this.resolved != nil {
		this.resolved.Write(to)
	} else {
		this.Await()
		this.resolved.Write(to)
	}
}

func (this *RESPFuture) Chan() chan<- RESPValue {
	return this.ch
}

func (this *RESPFuture) Await() RESPValue {
	if this.resolved != nil {
		return this.resolved
	} else {
		this.resolved = <-this.ch
		return this.resolved
	}
}
