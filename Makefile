build:
	gox -osarch="darwin/amd64" -output="./bin/pgweb_{{.OS}}_{{.Arch}}"

deps:
	go get github.com/mitchellh/gox