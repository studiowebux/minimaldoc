.PHONY: build build-server run run-server clean test docs install uninstall

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

# CLI
build:
	go build -o minimaldoc ./cmd/minimaldoc

run: build
	./minimaldoc build docs -o public

# Server (two ports: 8080 public, 8090 admin)
build-server:
	go build -o minimaldoc-server ./cmd/minimaldoc-server

run-server: build-server
	AUTH_JWT_SECRET=dev-secret-key-min-32-characters-long \
	DB_DRIVER=sqlite \
	DATABASE_URL=minimaldoc.db \
	SERVER_PORT=8080 \
	SERVER_ADMIN_PORT=8090 \
	./minimaldoc-server

run-server-pg: build-server
	AUTH_JWT_SECRET=dev-secret-key-min-32-characters-long \
	DB_DRIVER=postgres \
	DATABASE_URL="postgres://minimaldoc:minimaldoc@localhost:5432/minimaldoc?sslmode=disable" \
	SERVER_PORT=8080 \
	SERVER_ADMIN_PORT=8090 \
	./minimaldoc-server

# Bootstrap (create first site + admin) - uses admin port
bootstrap:
	@curl -s -X POST http://localhost:8090/api/bootstrap \
		-H "Content-Type: application/json" \
		-d '{"site_name":"My Docs","email":"admin@example.com"}' | jq .

# Docker
docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down

docker-logs:
	docker-compose -f docker/docker-compose.yml logs -f

# Common
clean:
	rm -rf minimaldoc minimaldoc-server public/ minimaldoc-server.db

test:
	go test ./...

docs: build
	./minimaldoc build docs -o public

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 minimaldoc $(DESTDIR)$(BINDIR)/minimaldoc

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/minimaldoc
