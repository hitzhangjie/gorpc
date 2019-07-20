PWD			:= $(shell pwd -P)
GOPATH		:= $(GOPATH):$(PWD)
BIN			:= "bin"

all:
	@GOPATH=$(GOPATH) go build src/test_nrpc.go
	@if [ ! -d $(BIN) ]; then mkdir $(BIN); fi
	@mv test_nrpc $(BIN)

.PHONY: clean
.PHONY: client

client:
	@GOPATH=$(GOPATH) go build client/test_nrpc_client.go
	@if [ ! -d $(BIN) ]; then mkdir $(BIN); fi
	@mv ./test_nrpc_client $(BIN)

clean:
	@rm -rf bin/test_nrpc
	@rm -rf bin/test_nrpc_client


