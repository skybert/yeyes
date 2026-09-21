BINARY := yeyes
PREFIX ?= $(HOME)/.local

SOURCES := $(wildcard *.go) $(wildcard *.m) go.mod go.sum

all: $(BINARY)

$(BINARY): $(SOURCES)
	mkdir -p build
	go build -o build/$(BINARY) .

run: $(BINARY)
	./build/$(BINARY)

install: $(BINARY)
	install -d $(PREFIX)/bin
	install -m 755 build/$(BINARY) $(PREFIX)/bin/$(BINARY)

clean:
	rm -rf build

.PHONY: all run install clean
