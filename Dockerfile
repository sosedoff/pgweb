FROM golang:1.3

COPY . /go/src/pgweb

WORKDIR /go/src/pgweb

RUN go get github.com/tools/godep
RUN	godep get github.com/jteeuwen/go-bindata/...
RUN godep restore

# need reconverts the static file
RUN	make build-asserts

RUN godep go build && \
    godep go install

ENTRYPOINT ["pgweb"]
