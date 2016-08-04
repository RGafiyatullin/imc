
export GOPATH = $(shell pwd)
export GOBIN = $(shell pwd)/bin

GO = go
RM_F = rm -f

IMCD = bin/imcd

all: $(IMCD)

clean: rm-bin-imcd

$(IMCD):
	$(GO) install src/github.com/rgafiyatullin/imc/server/imcd.go

rm-bin-imcd:
	$(RM_F) $(IMCD)




