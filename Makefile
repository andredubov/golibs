.PHONY:
.SILENT:

LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

get-deps:
	go get -u github.com/jackc/pgconn
	go get -u github.com/jackc/pgx/v4
	go get -u github.com/georgysavva/scany/pgxscan
	go get -u github.com/pkg/errors
	