# Boardsite Server
[![Go Report Card](https://goreportcard.com/badge/github.com/boardsite-io/server)](https://goreportcard.com/report/github.com/boardsite-io/server)

Websocket backend for [boardsite.io](https://boardsite.io).


### Build from source
* With Go tool (requires `go >=1.18` installed)
```
go build -o boardsite
```
* With docker (requires `docker` installed) 
```bash
docker build . --target bin --output .
```

The docker image containing the precompiled binary
can also be pulled from `ghcr.io/boardsite-io/boardsite-server:latest`


### Run Locally
Please use the provided docker-compose script to run a containerized 
instance together with the required redis cache locally.
```bash
make start
```
This will start the server on http://localhost:8000.

To stop the containers:
```
make stop
```

### Contribute

Contributions are always welcome. For small changes feel free to send us a PR. For bigger changes please create an issue
first to discuss your proposal.

### License
Licensed under [GNU AGPL v3.0](https://github.com/boardsite-io/boardsite-server/blob/master/LICENSE)