PWD			:= $(shell pwd -P)
GOPATH		:= $(GOPATH):$(PWD)
BIN			:= "bin"

all:
	@GOPATH=$(GOPATH) go build src/{{.ServerName}}.go
	@if [ ! -d $(BIN) ]; then mkdir $(BIN); fi
	@mv {{.ServerName}} $(BIN)

.PHONY: clean
.PHONY: client

client:
	@GOPATH=$(GOPATH) go build client/{{.ServerName}}_client.go
	@if [ ! -d $(BIN) ]; then mkdir $(BIN); fi
	@mv ./{{.ServerName}}_client $(BIN)

clean:
	@rm -rf bin/{{.ServerName}}
	@rm -rf bin/{{.ServerName}}_client


