package join

import "container/list"

type Awaitable interface {
	Await()
}

type Joinable interface {
	Join() Awaitable
}

func ReleaseJoiners(joiners *list.List) {
	for element := joiners.Front(); element != nil; element = element.Next() {
		element.Value.(chan<- bool) <- true
	}
}

func NewServerChan() chan chan<- bool {
	return make(chan chan<- bool, 32)
}

func NewChan() chan bool {
	return make(chan bool, 1)
}

func New(ch <-chan bool) Awaitable {
	a := new(awaitableAtChan)
	a.ch = ch
	return a
}

func NewStub() Awaitable {
	return new(awaitableStub)
}

type awaitableStub struct{}

func (this *awaitableStub) Await() {}

type awaitableAtChan struct {
	ch <-chan bool
}

func (this *awaitableAtChan) Await() {
	<-this.ch
}
