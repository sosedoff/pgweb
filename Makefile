TARGETS = darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
DOCKER_RELEASE_TAG = "sosedoff/pgweb:$(shell git describe --abbrev=0 --tags | sed 's/v//')"
DOCKER_LATEST_TAG = "sosedoff/pgweb:latest"
BINDATA_IGNORE = $(shell git ls-files -io --exclude-standard $< | sed 's/^/-ignore=/;s/[.]/[.]/g')

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
	go test -cover ./pkg/...

test-acceptance:
	@./script/test_ginkgo.sh

test-all:
	@./script/test_all.sh

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
		-ldflags "-X github.com/sosedoff/pgweb/pkg/command.GitCommit=$(GIT_COMMIT) -X github.com/sosedoff/pgweb/pkg/command.BuildTime=$(BUILD_TIME)" \
		-output "./bin/pgweb_{{.OS}}_{{.Arch}}"

	@echo "Building ARM binaries..."
	GOOS=linux GOARCH=arm GOARM=5 go build \
	  -ldflags "-X github.com/sosedoff/pgweb/pkg/command.GitCommit=$(GIT_COMMIT) -X github.com/sosedoff/pgweb/pkg/command.BuildTime=$(BUILD_TIME)" \
		-o "./bin/pgweb_linux_arm_v5"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

bootstrap:
	gox -build-toolchain

setup:
	go get github.com/tools/godep
	go get golang.org/x/tools/cmd/cover
	go get github.com/mitchellh/gox
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/sclevine/agouti
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	godep restore

clean:
	@rm -f ./pgweb
	@rm -rf ./bin/*
	@rm -f bindata.go

docker:
	docker build -t pgweb .

docker-release:
	docker build -t $(DOCKER_RELEASE_TAG) .
	docker build -t $(DOCKER_LATEST_TAG) .

docker-push:
	docker push $(DOCKER_RELEASE_TAG)
	docker push $(DOCKER_LATEST_TAG)
