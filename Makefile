.PHONY: all
all: build test

.PHONY: setup
setup:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HOME)/bin latest

.PHONY: build
build:
	go build -o jsonstore cmd/jsonstore/main.go

.PHONY: run
run:
	SERVE_PORT=8080 go run cmd/jsonstore/main.go

.PHONY: test
test:
	go test -count=1 -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html coverage.out -o coverage.html
	@coverage="$$(go tool cover -func coverage.out | grep 'total:' | awk '{print int($$3)}')"; \
	echo "The overall coverage is $$coverage%. Look at coverage.html for details.";

.PHONy: docker.build
docker.build:
	 docker build -t jsonstore -f build/Dockerfile .

.PHONY: fix
fix:
	$(GOPATH)/bin/golangci-lint run --fix

.PHONY: lint
lint:
	$(GOPATH)/bin/golangci-lint run

.PHONY: mocks.regenerate
mocks.regenerate:
	go get -u github.com/vektra/mockery/cmd/mockery
	mockery -dir=pkg/db -output=pkg/testlib/mocks -case underscore -all
	mockery -dir=pkg/service -output=pkg/testlib/mocks -case underscore -all
