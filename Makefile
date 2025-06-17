.DEFAULT_GOAL := test

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

staticcheck: fmt
	staticcheck ./...

revive: fmt
	revive -config revive.toml ./...

golangci: fmt
	golangci-lint run

lint: fmt staticcheck revive golangci

build: lint
	go build .

install: build
	go install .

test:
	go test -shuffle on .

testv:
	go test -shuffle on -v .

bench:
	go test -bench=. -benchmem -benchtime=5s -count=3 -run=NONE

.PHONY: fmt lint build install test testv bench
