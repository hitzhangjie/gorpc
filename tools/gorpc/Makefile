user := $(shell whoami)
inst := $(shell echo $(GOPATH) | cut -d':' -f1)

all: *.go
	@GOPATH=$(GOPATH) go build -o gorpc main.go

.PHONY: clean
.PHONY: install
.PHONY: uninstall

install:
	@mkdir -p $(inst)/bin/
	@cp ./gorpc $(inst)/bin/
ifeq ($(user),root)
#root, install for all user
	@[ -d /etc/gorpc ] || mkdir /etc/gorpc
	@cp -rf ./install/* /etc/gorpc/
else
#!root, install for current user
	@[ -d ~/.gorpc ] || mkdir ~/.gorpc
	@cp -rf ./install/* ~/.gorpc/
endif
	@echo "install finished"

uninstall:
	@rm -rf $(inst)/bin/gorpc
ifeq ($(user), root)
#root, install for all user
	@rm -rf /etc/gorpc
else
#!root, install for current user
	@rm -rf ~/.gorpc
endif
	@echo "uninstall finished"

clean:
	@rm ./gorpc -f
	@echo "clean finished"

