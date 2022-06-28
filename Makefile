TARGETS = darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION = $(shell go version | awk {'print $$3'})
DOCKER_RELEASE_TAG = "sosedoff/pgweb:$(shell git describe --abbrev=0 --tags | sed 's/v//')"
DOCKER_LATEST_TAG = "sosedoff/pgweb:latest"
LDFLAGS = -s -w
PKG = github.com/sosedoff/pgweb

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make setup           : Install all necessary dependencies"
	@echo "make dev             : Generate development build"
	@echo "make build           : Generate production build for current OS"
	@echo "make bootstrap       : Install cross-compilation toolchain"
	@echo "make release         : Generate binaries for all supported OSes"
	@echo "make test            : Execute test suite"
	@echo "make test-all        : Execute test suite on multiple PG versions"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make docker          : Build docker image"
	@echo "make docker-release  : Build and tag docker image"
	@echo "make docker-push     : Push docker images to registry"
	@echo ""

test:
	go test -race -cover ./pkg/...

test-all:
	@./script/test_all.sh
	@./script/test_cockroach.sh

dev:
	go build
	@echo "You can now execute ./pgweb"

build:
	go build
	@echo "You can now execute ./pgweb"

release: LDFLAGS += -X $(PKG)/pkg/command.GitCommit=$(GIT_COMMIT)
release: LDFLAGS += -X $(PKG)/pkg/command.BuildTime=$(BUILD_TIME)
release: LDFLAGS += -X $(PKG)/pkg/command.GoVersion=$(GO_VERSION)
release:
	@echo "Building binaries..."
	@gox \
		-osarch "$(TARGETS)" \
		-ldflags "$(LDFLAGS)" \
		-output "./bin/pgweb_{{.OS}}_{{.Arch}}"

	@echo "Building ARM binaries..."
	GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o "./bin/pgweb_linux_arm_v5"

	@echo "Building ARM64 binaries..."
	GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags "$(LDFLAGS)" -o "./bin/pgweb_linux_arm64_v7"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

bootstrap:
	gox -build-toolchain

setup:
	go install github.com/mitchellh/gox@v1.0.1

clean:
	@rm -f ./pgweb
	@rm -rf ./bin/*

docker:
	docker build --no-cache -t pgweb .

docker-release:
	docker build --no-cache -t $(DOCKER_RELEASE_TAG) .
	docker tag $(DOCKER_RELEASE_TAG) $(DOCKER_LATEST_TAG)
	docker images $(DOCKER_RELEASE_TAG)

docker-push:
	docker push $(DOCKER_RELEASE_TAG)
	docker push $(DOCKER_LATEST_TAG)
