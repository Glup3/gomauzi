# GoMauzi
API server for Mauzi - spending tracker

## Running the server
The Pokemon dex number of Mauzi is 52. The default server port is `8052`

```
go run main.go
```

## Local Database with Docker

```
docker compose up -d
```

## Running Tests
You need to have Docker installed and running.

```
go test -v ./...
```