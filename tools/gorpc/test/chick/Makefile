all: test_chick.proto
	gorpc create -protodir=. -protofile=test_chick.proto -protocol=chick -v

.PHONY: clean

clean:
	rm -rf fruit_svr
