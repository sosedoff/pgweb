FROM golang:1.3

COPY . /go/src/pgweb
WORKDIR /go/src/pgweb

RUN go get github.com/tools/godep

RUN godep get github.com/mitchellh/gox
RUN	godep get github.com/jteeuwen/go-bindata/...

RUN godep restore
RUN	go-bindata -ignore=\\.gitignore -ignore=\\.DS_Store -ignore=\\.gitkeep static/...

RUN godep go build && \
    godep go install

ENTRYPOINT ["pgweb"]
