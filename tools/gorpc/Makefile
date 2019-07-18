user := $(shell whoami)

all: *.go
	@GOPATH=$(GOPATH) go build -o goneat main.go

.PHONY: clean
.PHONY: install

install:
ifeq ($(user),root)
#root, install for all user
	@cp ./goneat /usr/bin/
	@ln -sf /usr/bin/goneat /usr/bin/gogen
	@[ -d /etc/goneat ] || mkdir /etc/goneat
	@cp -rf ./asset /etc/goneat/
	@cp -rf ./asset_taf /etc/goneat/
	@cp -f ./spec/protospec.json /etc/goneat/
	@cp -f ./spec/deployspec.json /etc/goneat/
else
#!root, install for current user
	@mkdir -p ~/bin/
	@cp ./goneat ~/bin/
	@ln -sf ~/bin/goneat ~/bin/gogen
	@[ -d ~/.goneat ] || mkdir ~/.goneat
	@cp -rf ./asset ~/.goneat/
	@cp -rf ./asset_taf ~/.goneat/
	@cp -f ./spec/protospec.json ~/.goneat/
	@cp -f ./spec/deployspec.json ~/.goneat/
endif
	@echo "install finished"
clean:
	@rm ./goneat -f
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

