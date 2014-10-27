dev:
	rm -f bindata.go
	go-bindata -debug -ignore=\\.gitignore -ignore=\\.DS_Store -ignore=\\.gitkeep static/...
	go build
	@echo "You can now execute ./pgweb"

build:
	rm -f bindata.go
	go-bindata -ignore=\\.gitignore -ignore=\\.DS_Store -ignore=\\.gitkeep static/...
	gox -osarch="darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386" -output="./bin/pgweb_{{.OS}}_{{.Arch}}"

setup:
	go get github.com/mitchellh/gox
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -debug -ignore=\\.gitignore -ignore=\\.DS_Store -ignore=\\.gitkeep static/...
	go get

clean:
	rm -f ./pgweb
	rm -f ./bin/*
	rm -f bindata.go