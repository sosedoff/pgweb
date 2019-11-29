TARGETS = darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION = $(shell go version | awk {'print $$3'})
BINDATA_IGNORE = $(shell git ls-files -io --exclude-standard $< | sed 's/^/-ignore=/;s/[.]/[.]/g')
DOCKER_RELEASE_TAG = "sosedoff/pgweb:$(shell git describe --abbrev=0 --tags | sed 's/v//')"
DOCKER_LATEST_TAG = "sosedoff/pgweb:latest"

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
	@echo "make assets          : Generate production assets file"
	@echo "make dev-assets      : Generate development assets file"
	@echo "make docker          : Build docker image"
	@echo "make docker-release  : Build and tag docker image"
	@echo "make docker-push     : Push docker images to registry"
	@echo ""

test:
	go test -race -cover ./pkg/...

test-all:
	@./script/test_all.sh
	@./script/test_cockroach.sh

assets: static/
	go-bindata -o pkg/data/bindata.go -pkg data $(BINDATA_OPTS) $(BINDATA_IGNORE) -ignore=[.]gitignore -ignore=[.]gitkeep $<...

dev-assets:
	@$(MAKE) --no-print-directory assets BINDATA_OPTS="-debug"

dev: dev-assets
	go build
	@echo "You can now execute ./pgweb"

build: assets
	go build
	@echo "You can now execute ./pgweb"

release: clean assets
	@echo "Building binaries..."
	@gox \
		-osarch "$(TARGETS)" \
		-ldflags "-X github.com/sosedoff/pgweb/pkg/command.GitCommit=$(GIT_COMMIT) -X github.com/sosedoff/pgweb/pkg/command.BuildTime=$(BUILD_TIME) -X github.com/sosedoff/pgweb/pkg/command.GoVersion=$(GO_VERSION)" \
		-output "./bin/pgweb_{{.OS}}_{{.Arch}}"

	@echo "Building ARM binaries..."
	GOOS=linux GOARCH=arm GOARM=5 go build \
	  -ldflags "-X github.com/sosedoff/pgweb/pkg/command.GitCommit=$(GIT_COMMIT) -X github.com/sosedoff/pgweb/pkg/command.BuildTime=$(BUILD_TIME) -X github.com/sosedoff/pgweb/pkg/command.GoVersion=$(GO_VERSION)" \
		-o "./bin/pgweb_linux_arm_v5"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

bootstrap:
	gox -build-toolchain

setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/mitchellh/gox
	go get -u github.com/go-bindata/go-bindata/...
	dep ensure

clean:
	@rm -f ./pgweb
	@rm -rf ./bin/*
	@rm -f bindata.go

docker:
	docker build -t pgweb .

docker-release:
	docker build --no-cache -t $(DOCKER_RELEASE_TAG) .
	docker tag $(DOCKER_RELEASE_TAG) $(DOCKER_LATEST_TAG)
	docker images $(DOCKER_RELEASE_TAG)

docker-push:
	docker push $(DOCKER_RELEASE_TAG)
	docker push $(DOCKER_LATEST_TAG)
