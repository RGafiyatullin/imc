package join

// Provides an interface to "join" an actor (implented by actor-handle)
type Joinable interface {
	// does not block itself instead returning kind of Future one can wait for.
	Join() Awaitable
}
