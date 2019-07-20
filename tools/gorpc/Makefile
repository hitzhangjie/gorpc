user := $(shell whoami)

all: *.go
	@GOPATH=$(GOPATH) go build -o gorpc main.go

.PHONY: clean
.PHONY: install

install:
ifeq ($(user),root)
#root, install for all user
	@cp ./gorpc /usr/bin/
	@[ -d /etc/gorpc ] || mkdir /etc/gorpc
	@cp -rf ./asset /etc/gorpc/
	@cp -f ./spec/protospec.json /etc/gorpc/
else
#!root, install for current user
	@mkdir -p ~/bin/
	@cp ./gorpc ~/bin/
	@[ -d ~/.gorpc ] || mkdir ~/.gorpc
	@cp -rf ./asset ~/.gorpc/
	@cp -f ./spec/protospec.json ~/.gorpc/
endif
	@echo "install finished"

clean:
	@rm ./gorpc -f
	@rm ./test_ilive -rf
	@rm ./test_nrpc -rf
	@rm ./test_simplesso -rf
	@echo "clean finished"

test_ilive: 
	@go run main.go -protodir=test/ilive -protofile=test_ilive.proto -protocol=ilive -verbose=false
	@echo "test_ilive generated"

test_nrpc:
	@go run main.go -protodir=test/nrpc -protofile=test_nrpc.proto -protocol=nrpc
	@echo "test_nrpc generated"

test_simplesso:
	@go run main.go -protodir=test/simplesso -protofile=test_simplesso.proto -protocol=simplesso
	@echo "test_simplesso generated"

