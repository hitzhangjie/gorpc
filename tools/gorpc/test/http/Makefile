all: test_gorpc.proto
	goneat create -protodir=. -protofile=test_gorpc.proto -protocol=gorpc -httpon

.PHONY: clean

clean:
	rm -rf test_gorpc
