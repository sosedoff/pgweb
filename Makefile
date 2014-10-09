build:
	gox -osarch="linux/amd64 darwin/amd64" -output="./bin/pgweb_{{.OS}}_{{.Arch}}"

deps:
	go get github.com/mitchellh/gox