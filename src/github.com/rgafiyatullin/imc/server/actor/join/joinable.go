package join

type Joinable interface {
	Join() Awaitable
}
