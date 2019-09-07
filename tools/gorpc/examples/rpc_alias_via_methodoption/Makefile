all:
	@echo "------------- [gorpc create] --------------"
	gorpc create -protofile=helloworld.proto -protocol=whisper

.PHONY: clean

clean:
	rm -rf greeter

compile:
	@echo "------------- [compile gorpc] -------------"
	cd ../../../gorpc/
	make clean && make && make install
	@cd -
