.PHONY: test test-coverage lint fmt

all:

test:
	go test $(COVERFLAG) -v ./...

test-coverage: COVERFLAG=-cover
test-coverage: test

lint:
	golangci-lint run

fmt:
	go fmt \
		github.com/scality/bucketclient/go
