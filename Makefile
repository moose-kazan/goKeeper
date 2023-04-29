.PHONY: all binary

PREFIX=/usr/local/

all: binary

binary:
	mkdir -p ${PWD}/build/bin
	go build -o ${PWD}/build/bin/goKeeperViewer ${PWD}/cmd/main.go

clean:
	rm -rf build

install: all
	mkdir -p ${PREFIX}/bin
	cp -v ${PWD}/build/bin/goKeeperViewer ${PREFIX}/bin/goKeeperViewer

