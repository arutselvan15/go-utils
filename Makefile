APP_NAME=utils
APP_DISPLAY_NAME=Utils
PROJECT_NAME=go-${APP_NAME}

GOLINTV=v1.16.0
GOARGS=GOOS=linux GOARCH=amd64 GO111MODULE=on
GO=${GOARGS} go
GOLINT=${GOARGS} golangci-lint

BUILD=$(or ${BUILD_NUMBER},0)
VERSION=v0.1.${BUILD}
DATE=$(shell date)
HOSTNAME=$(shell hostname)

DOCKER_HUB=arutselvan15

## display this help message
help:
	@echo ''
	@echo 'make targets for the project:'
	@echo
	@echo 'Usage:'
	@echo '  ## Develop / Test Commands'
	@echo '  all             Run all the commands.'
	@echo '  clean           Run clean up.'
	@echo '  deps            Run all dependencies that are needed.'
	@echo '  fmt             Run code formatter.'
	@echo '  check           Run static code analysis (lint).'
	@echo '  test            Run tests on project.'
	@echo '  build           Run build to build the project.'
	@echo '  run             Run run on run the binary.'

	@echo '  ### one time setup ###'
	@echo '  setup           Run setup of tools needed.'
	@echo '  download-go     Run download-go to download go.'
	@echo '  install-go      Run install-go to install go.'
	@echo '  mod-init        Run mod-init to init go mods setup. This should only be executed once.'
	@echo '  install-golint  Run install-golint to install go lint.'

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

build:
	@echo "==> Build local..."
	echo "Version=${VERSION}" > version.txt
	echo "Date=${DATE}" >> version.txt
	echo "Host=${HOSTNAME}" >> version.txt
	cat version.txt
	CGO_ENABLED=0 ${GO} build -o bin/${PROJECT_NAME} .

run:
	@echo "==> Run local..."
	cd bin;./${PROJECT_NAME}

test: clean
	@echo "==> Testing..."
	CGO_ENABLED=0 ${GO} test -v -covermode=atomic -count=1 ./... -coverprofile coverage.out
	${GO} test -race -covermode=atomic -count=1 ./... -json > report.json
	${GO} tool cover -func=coverage.out

### one time setup ###
setup: mod-init install-golint

mod-init:
	@echo "==> Mod Init..."
	${GO} mod init

install-golint:
	@echo "==> Install go lint..."
	${GO} get github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLINTV}
