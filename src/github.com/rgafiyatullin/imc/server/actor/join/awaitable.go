package join

// The Future to wait for (see Joinable)
type Awaitable interface {
	// Blocks the caller until signaled by the awaited party
	Await()
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
