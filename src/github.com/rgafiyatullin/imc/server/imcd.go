package main

import (
	"container/list"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/actor/join"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/netsrv"
	"github.com/rgafiyatullin/imc/server/storage"
	"runtime"
)

func main() {
	awaitList := list.New()
	defer awaitBeforeExit(awaitList)

	config := config.New()

	topActorCtx := actor.New(config)
	logctx := topActorCtx.Log()
	logctx.Info("imcd: starting up")
	logctx.Info("runtime.NumCPU: %v", runtime.NumCPU())
	logctx.Info("runtime.GOMAXPROCS: %v", runtime.GOMAXPROCS(0))

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
