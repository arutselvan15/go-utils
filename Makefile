GO=GOOS=linux GOARCH=amd64 GO111MODULE=on go
APPNAME=goUtils
BINARY=go-utils
PORT=5000
GOVERSION=go1.11.5
GOLINTV=v1.16.0
VERSION=v1.0.0

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
	@echo '  gen-server      Run gen-server to generate swagger server.'
	@echo '  push            Run push to push changes to git.'
	@echo '  build           Run build to build the project.'
	@echo '  run             Run run on run the binary.'
	@echo '  docker-build    Run docker-build to build the docker image.'
	@echo '  docker-run      Run docker-run on run the app in docker container.'

	@echo '  ### one time setup ###'
	@echo '  setup           Run setup of tools needed.'
	@echo '  download-go     Run download-go to download go.'
	@echo '  install-go      Run install-go to install go.'
	@echo '  mod-init        Run mod-init to init go mods setup. This should only be executed once.'
	@echo '  install-swagger Run install-swagger to install swagger.'
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
	golangci-lint run --fast --tests

gen-server:
	@echo "==> Generate Swagger Server..."
	swagger generate server -f api.yaml -A ${APPNAME}

delete-server:
	@echo "==> Deleting Swagger Server..."
	rm -rf cmd models restapi

gen-client:
	@echo "==> Generate Swagger Client..."
	swagger generate client -f api.yaml -A ${APPNAME}

build:
	@echo "==> Build local..."
	echo ${VERSION} > version.txt
	CGO_ENABLED=0 ${GO} build -o bin/${BINARY} ./cmd/${BINARY}/

run:
	@echo "==> Run local..."
	cd bin;./${BINARY} --port=${PORT}

docker-build:
	@echo "==> Build Docker Image..."
	docker build -t ${BINARY}:${VERSION} .

docker-run:
	@echo "==> Docker Run..."
	docker run -p ${PORT}:${PORT} ${BINARY}:${VERSION}

test: clean
	@echo "==> Testing..."
	CGO_ENABLED=0 ${GO} test -v -covermode=atomic -count=1 ./... -coverprofile coverage.out
	${GO} test -race -covermode=atomic -count=1 ./... -json > report.json
	${GO} tool cover -func=coverage.out

push:
	@echo "==> Git Push..."
	git push

### one time setup ###
setup: mod-init install-golint install-swagger

download-go:
	cd ${HOME}/Downloads;wget https://dl.google.com/go/${GOVERSION}.linux-amd64.tar.gz

install-go:
	cd ${HOME}/Downloads;chmod +x ${GOVERSION}.linux-amd64.tar.gz; tar -xvf ${GOVERSION}.linux-amd64.tar.gz
	rm -rf /usr/local/go
	mv ${HOME}/Downloads/go /usr/local

mod-init:
	@echo "==> Mod Init..."
	${GO} mod init

install-swagger: swagger-deps
	@echo "==> Install go swagger..."
	${GO} get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger-deps:
	@echo "==> Install swagger dependencies..."
	${GO} get -u github.com/go-openapi/runtime
	${GO} get -u github.com/jessevdk/go-flags

install-golint:
	@echo "==> Install go lint..."
	${GO} get github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLINTV}
