all: test_swan.proto
	gorpc create -protodir=. -protofile=test_swan.proto -protocol=swan -v

.PHONY: clean

clean:
	rm -rf fruit_svr
