
GO = go

export GOPATH = $(shell pwd)

export GOBIN = $(shell pwd)/bin

all: bin/imcd

bin/imcd:
	$(GO) install src/github.com/rgafiyatullin/imc/server/imcd.go



