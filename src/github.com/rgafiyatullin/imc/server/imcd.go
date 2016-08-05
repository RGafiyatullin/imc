package main

import (
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/config"
	"github.com/rgafiyatullin/imc/server/netsrv"
	"github.com/rgafiyatullin/imc/server/storage"
)

func main() {
	fmt.Println("Helloes! I'm the IMC daemon")

	topActorCtx := actor.NewCtx()
	topActorCtx.Log().Info("System start")

	config := config.New()

	storageSup := storage.StartSup(topActorCtx.NewChild("storage_sup"), config)
	ringmgr := storageSup.QueryRingMgr()
	listener, _ := netsrv.StartListener(topActorCtx.NewChild("listener"), config, ringmgr)

	joinStorage := storageSup.Join()
	joinListener := listener.Join()

	joinStorage.Await()
	joinListener.Await()
}
