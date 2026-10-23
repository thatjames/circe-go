# Build entry point for circe; CI calls it. See features/build/makefile.md.

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
IMAGE   ?= circe
PLATFORMS ?= linux/amd64,linux/arm64
LDFLAGS := -s -w -X github.com/thatjames/circe-go/internal/cli.version=$(VERSION)

.PHONY: build docker docker-push dev dev-down clean

build:
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o circe ./cmd/circe

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
