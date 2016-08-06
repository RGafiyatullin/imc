package main

import (
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/netsrv"
	"github.com/rgafiyatullin/imc/server/storage"
	"container/list"
	"github.com/rgafiyatullin/imc/server/actor/join"
)

func main() {
	fmt.Println("Helloes! I'm the IMC daemon")

	awaitList := list.New()
	defer awaitBeforeExit(awaitList)

	config := config.New()

	topActorCtx := actor.New(config)
	topActorCtx.Log().Info("starting up")

	storageSup := storage.StartSup(topActorCtx.NewChild("storage_sup"), config)
	awaitList.PushBack(storageSup.Join())

	ringmgr := storageSup.QueryRingMgr()
	listener, listenErr := netsrv.StartListener(topActorCtx.NewChild("listener"), config, ringmgr)
	if listenErr == nil {
		awaitList.PushBack(listener.Join())
	}
}

func awaitBeforeExit(awaitList *list.List) {
	for elt := awaitList.Front(); elt != nil; elt = elt.Next() {
		elt.Value.(join.Awaitable).Await()
	}
}
