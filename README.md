# Recorder

## Development

### Generate the protobuf code

```shell
brew install protobuf
protoc -Iprotos/v1 --go_out=. --go-grpc_out=. ./protos/v1/api.proto
```

### Add a SQL migration

```shell
brew install golang-migrate
migrate create -dir ./sql -ext .sql <NAME>
```

### Lint the code

```shell
brew install golangci-lint hadolint sqlfluff
golangci-lint run
hadolint ./images/*.Dockerfile
sqlfluff lint --dialect postgres ./sql/*.sql
```

### Run unit tests

```shell
go test -v ./...
```

### Generate certs

```shell
brew install mkcert
mkcert -install
mkdir -p ./certs
mkcert -cert-file ./certs/postgresql.crt.pem -key-file ./certs/postgresql.key.pem postgresql localhost
mkcert -cert-file ./certs/server.crt.pem -key-file ./certs/server.key.pem server localhost
export CAROOT=$(mkcert -CAROOT)
```

### Run integration tests

```shell
docker compose -f ./deploy/dev/docker-compose.yaml up -d --build
docker compose -f ./deploy/dev/docker-compose.yaml stop
docker compose -f ./deploy/dev/docker-compose.yaml rm
```

### Run locally

```shell
docker compose -f ./deploy/local/docker-compose.yaml up -d --build
docker compose -f ./deploy/local/docker-compose.yaml stop
docker compose -f ./deploy/local/docker-compose.yaml rm
go run ./cmd/cli version
```
