PKG = github.com/sosedoff/pgweb
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION ?= $(shell go version | awk {'print $$3'})

DOCKER_RELEASE_TAG = "sosedoff/pgweb:$(shell git describe --abbrev=0 --tags | sed 's/v//')"
DOCKER_LATEST_TAG = "sosedoff/pgweb:latest"

LDFLAGS = -s -w
LDFLAGS += -X $(PKG)/pkg/command.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X $(PKG)/pkg/command.BuildTime=$(BUILD_TIME)
LDFLAGS += -X $(PKG)/pkg/command.GoVersion=$(GO_VERSION)

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make dev             : Generate development build"
	@echo "make build           : Generate production build for current OS"
	@echo "make release         : Generate binaries for all supported OSes"
	@echo "make test            : Execute test suite"
	@echo "make test-all        : Execute test suite on multiple PG versions"
	@echo "make lint            : Execute code linter"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make docker          : Build docker image"
	@echo "make docker-release  : Build and tag docker image"
	@echo "make docker-push     : Push docker images to registry"
	@echo ""

test:
	go test -v -race -cover ./pkg/...

test-all:
	@./script/test_all.sh
	@./script/test_cockroach.sh

lint:
	golangci-lint run

dev:
	go build
	@echo "You can now execute ./pgweb"

build:
	go build -ldflags '${LDFLAGS}'
	@echo "You can now execute ./pgweb"

install:
	go install -ldflags '${LDFLAGS}'
	@echo "You can now execute pgweb"

release: clean
	@echo "Building binaries..."
	@LDFLAGS='${LDFLAGS}' ./script/build_all.sh

clean:
	@echo "Removing all artifacts"
	@rm -rf ./pgweb ./bin/*

docker:
	docker build --no-cache -t pgweb .

docker-run:
	docker run --rm -p 8081:8081 -it pgweb

docker-release:
	docker build --no-cache -t $(DOCKER_RELEASE_TAG) .
	docker tag $(DOCKER_RELEASE_TAG) $(DOCKER_LATEST_TAG)
	docker images $(DOCKER_RELEASE_TAG)

docker-push:
	docker push $(DOCKER_RELEASE_TAG)
	docker push $(DOCKER_LATEST_TAG)
