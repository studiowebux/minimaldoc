.PHONY: build run clean test docs install uninstall

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

build:
	go build -o minimaldoc ./cmd/minimaldoc

run: build
	./minimaldoc build docs -o public

clean:
	rm -rf minimaldoc public/

test:
	go test ./...

docs: build
	./minimaldoc build docs -o public

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 minimaldoc $(DESTDIR)$(BINDIR)/minimaldoc

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/minimaldoc
