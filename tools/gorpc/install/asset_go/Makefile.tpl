{{- $svrName := (index .Services 0).Name -}}
PWD			:= $(shell pwd -P)
GOPATH		:= $(GOPATH):$(PWD)
BIN			:= "bin"

all:
	@GOPATH=$(GOPATH) go build src/{{$svrName}}.go
	@if [ ! -d $(BIN) ]; then mkdir $(BIN); fi
	@mv {{$svrName}} $(BIN)

.PHONY: clean
.PHONY: client

client:
	@GOPATH=$(GOPATH) go build client/{{$svrName}}_client.go
	@if [ ! -d $(BIN) ]; then mkdir $(BIN); fi
	@mv ./{{$svrName}}_client $(BIN)

clean:
	@rm -rf bin/{{$svrName}}
	@rm -rf bin/{{$svrName}}_client


