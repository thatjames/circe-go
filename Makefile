# Build entry point for circe; CI calls it. See features/build/makefile.md.

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
IMAGE   ?= circe
PLATFORMS ?= linux/amd64,linux/arm64
LDFLAGS := -s -w -X github.com/thatjames/circe-go/internal/cli.version=$(VERSION)
# The one pin lives in the Lint workflow; read it from there.
GOLANGCI_LINT_VERSION ?= $(shell sed -n 's/^ *version: *//p' .github/workflows/lint.yml)

.PHONY: build docker docker-push dev dev-down clean check lint vuln

build:
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o circe ./cmd/circe

# Every check CI runs, in CI's order; stops at the first failure.
check:
	go mod tidy -diff
	CGO_ENABLED=0 go build ./...
	$(MAKE) --no-print-directory build
	CGO_ENABLED=0 go vet ./...
	CGO_ENABLED=0 go test ./...
	$(MAKE) --no-print-directory lint
	$(MAKE) --no-print-directory vuln

lint:
	CGO_ENABLED=0 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

vuln:
	CGO_ENABLED=0 go run golang.org/x/vuln/cmd/govulncheck@latest ./...

docker:
	docker build --load --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) .

docker-push:
	docker buildx build --push --platform $(PLATFORMS) --build-arg VERSION=$(VERSION) \
		-t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

dev:
	docker compose up -d --wait

dev-down:
	docker compose down

clean:
	rm -rf circe dist
