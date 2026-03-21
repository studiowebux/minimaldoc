.PHONY: build clean test lint ci docs check

# Build the minimaldoc CLI binary
build:
	go build -o minimaldoc ./cmd/minimaldoc

# Build docs site (docs/ → public/)
docs: build
	./minimaldoc build docs -o public

# Run tests
test:
	go test ./...

# Lint (gofmt + go vet)
lint:
	@echo "Running gofmt..."
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	@echo "Running go vet..."
	go vet ./...

# CI pipeline (lint + test)
ci: lint test

# Static analysis
check:
	@which staticcheck > /dev/null 2>&1 || (echo "Install: go install honnef.co/go/tools/cmd/staticcheck@latest" && exit 1)
	staticcheck ./...

# Clean build artifacts
clean:
	rm -f minimaldoc
	rm -rf public/
