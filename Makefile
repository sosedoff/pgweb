build:
	rm -f bindata.go
	go-bindata -ignore=\\.gitignore -ignore=\\.DS_Store -ignore=\\.gitkeep static/...
	gox -osarch="darwin/amd64 linux/amd64 linux/386 windows/amd64 windows/386" -output="./bin/pgweb_{{.OS}}_{{.Arch}}"

dev:
	rm -f bindata.go
	go-bindata -debug -ignore=\\.gitignore -ignore=\\.DS_Store -ignore=\\.gitkeep static/...
	go build
	@echo "You can now execute ./pgweb"

deps:
	go get
	go get github.com/mitchellh/gox
	go get github.com/jteeuwen/go-bindata/...

clean:
	rm -f ./pgweb
	rm -f ./bin/*
	rm -r bindata.go