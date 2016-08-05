package join

import "container/list"

func ReleaseJoiners(joiners *list.List) {
	for element := joiners.Front(); element != nil; element = element.Next() {
		element.Value.(chan<- bool) <- true
	}
}

func NewServerChan() chan chan<- bool {
	return make(chan chan<- bool, 32)
}

func NewClientChan() chan bool {
	return make(chan bool, 1)
}
