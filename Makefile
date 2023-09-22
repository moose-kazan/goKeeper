.PHONY: all binary

PREFIX ?= /usr

GOPATH=$(shell go env GOPATH)

all: binary

binary:
	go build -o goKeeperViewer

install: all
	#$(GOPATH)/bin/fyne install
	mkdir -p $(PREFIX)/bin
	install -m 755 ./goKeeperViewer $(PREFIX)/bin/

	mkdir -p $(PREFIX)/share/applications
	install -m 644 ./goKeeperViewer.desktop $(PREFIX)/share/applications/
	
	mkdir -p $(PREFIX)/share/pixmaps
	install -m 644 ./Icon.png $(PREFIX)/share/pixmaps/goKeeperViewer.png

clean:
	rm -f goKeeperViewer
