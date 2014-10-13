build:
	gox -osarch="darwin/amd64" -output="./bin/pgweb_{{.OS}}_{{.Arch}}"

deps:
	go get github.com/mitchellh/gox
	go get github.com/jteeuwen/go-bindata/...

clean:
	rm -f ./pgweb
	rm -f ./bin/*