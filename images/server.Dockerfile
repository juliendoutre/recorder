# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1
ARG ALPINE_VERSION=3.19
ARG DEBIAN_VERSION=12

FROM --platform=$BUILDPLATFORM index.docker.io/golang:$GO_VERSION-alpine$ALPINE_VERSION AS builder

ARG TARGETOS
ARG TARGETARCH
ARG GO_VERSION

WORKDIR /recorder/server

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/server ./cmd/server
COPY ./pkg ./pkg
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-s -w -X main.GoVersion=$GO_VERSION -X main.Os=$TARGETOS -X main.Arch=$TARGETARCH" -o /server ./cmd/server

FROM --platform=$TARGETPLATFORM gcr.io/distroless/base-debian$DEBIAN_VERSION:latest AS runner

LABEL org.opencontainers.image.authors Julien Doutre <jul.doutre@gmail.com>
LABEL org.opencontainers.image.title recorder
LABEL org.opencontainers.image.url https://github.com/juliendoutre/recorder
LABEL org.opencontainers.image.documentation https://github.com/juliendoutre/recorder
LABEL org.opencontainers.image.source https://github.com/juliendoutre/recorder
LABEL org.opencontainers.image.licenses MIT

WORKDIR /

COPY --from=builder /server /server

USER nonroot:nonroot

EXPOSE 8000

ENTRYPOINT ["/server"]
