.PHONY: all binary

PREFIX=/usr/local/

all: binary

binary:
	fyne build

install: all
	fyne install
