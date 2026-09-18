BINARY := geyes
PREFIX ?= $(HOME)/.local

SOURCES := $(wildcard *.go) $(wildcard *.m) go.mod go.sum

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY) .

run: $(BINARY)
	./$(BINARY)

install: $(BINARY)
	install -d $(PREFIX)/bin
	install -m 755 $(BINARY) $(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY)

.PHONY: all run install clean
