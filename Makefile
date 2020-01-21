GO=GO111MODULE=on go
GO_CROSS_CMPL=GOOS=linux GOARCH=amd64 ${GO}

BINARY=bin/go-utils
MAIN_GO=cmd/main.go

BUILD=$(or ${BUILD_NUMBER},0)
VERSION=v1.0.${BUILD}
DATE=$(shell date)
HOSTNAME=$(shell hostname)

DOCKER_HUB=arutselvan15

all: clean deps fmt check test

clean:
	@echo "==> Cleaning..."
	rm -f report.json coverage.out

deps:
	@echo "==> Getting Dependencies..."
	${GO} mod tidy
	${GO} mod download

fmt:
	@echo "==> Code Formatting..."
	${GO} fmt ./...

check: fmt
	@echo "==> Code Check..."
	${GOLINT} run --fast --tests

build: test gen-version
	@echo "==> Build Local..."
	CGO_ENABLED=0 ${GO} build -o ${BINARY} ${MAIN_GO}

gen-version:
	@echo "==> Generating Version..."
	echo "Version=${VERSION}" > version.txt
	echo "Date=${DATE}" >> version.txt
	echo "Host=${HOSTNAME}" >> version.txt
	cat version.txt

test: clean
	@echo "==> Testing..."
	CGO_ENABLED=0 ${GO} test -v -covermode=atomic -count=1 ./... -coverprofile coverage.out
	CGO_ENABLED=1 ${GO} test -race -covermode=atomic -count=1 ./... -json > report.json
	${GO} tool cover -func=coverage.out