all: test_gorpc.proto
	gorpc create -protodir=. -protofile=test_gorpc.proto -protocol=gorpc -v -g

.PHONY: clean

clean:
	rm -rf fruit_svr
