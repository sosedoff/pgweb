BINDATA_IGNORE = $(shell git ls-files -io --exclude-standard $< | sed 's/^/-ignore=/;s/[.]/[.]/g')

dev: dev-assets
	godep go build
	@echo "You can now execute ./pgweb"

assets: static/
	go-bindata $(BINDATA_OPTS) $(BINDATA_IGNORE) -ignore=[.]gitignore -ignore=[.]gitkeep $<...

dev-assets:
	@$(MAKE) --no-print-directory assets BINDATA_OPTS="-debug"

build: assets
	gox -osarch="darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386" -output="./bin/pgweb_{{.OS}}_{{.Arch}}"

setup:
	go get github.com/tools/godep
	godep get github.com/mitchellh/gox
	godep get github.com/jteeuwen/go-bindata/...
	godep restore

clean:
	rm -f ./pgweb
	rm -f ./bin/*
	rm -f bindata.go

docker:
	docker build -t pgweb .