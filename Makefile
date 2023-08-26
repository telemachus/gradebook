.DEFAULT_GOAL := lint

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

vet: fmt
	go vet .

lint: vet
	golangci-lint run

build: lint
	go build .

install: build
	go install .

test:
	go test -shuffle on .

testv:
	go test -shuffle on -v .

clean:
	go clean -i -r -cache

.PHONY: fmt vet lint build install test testv clean
