# syntax=docker/dockerfile:1
# See features/build/docker-image.md.

FROM --platform=$BUILDPLATFORM golang:1.27.2 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath \
      -ldflags "-s -w -X github.com/thatjames/circe-go/internal/cli.version=${VERSION}" \
      -o /out/circe ./cmd/circe

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.source=https://github.com/thatjames/circe-go
COPY --from=build /out/circe /usr/local/bin/circe
ENTRYPOINT ["circe"]
