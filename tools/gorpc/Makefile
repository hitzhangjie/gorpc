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
else
#!root, install for current user
	@mkdir -p ~/bin/
	@cp ./gorpc ~/bin/
	@[ -d ~/.gorpc ] || mkdir ~/.gorpc
	@cp -rf ./asset ~/.gorpc/
endif
	@echo "install finished"

clean:
	@rm ./gorpc -f
	@echo "clean finished"

