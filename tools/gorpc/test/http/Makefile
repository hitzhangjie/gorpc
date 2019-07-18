all: test_nrpc.proto
	#protoc --neat_out=plugins=nserver:. Test.proto
	goneat test_nrpc.proto nrpc httpon

.PHONY: clean

clean:
	rm -rf test_nrpc
