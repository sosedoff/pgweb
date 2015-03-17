FROM golang:1.4.2

COPY . /go/src/pgweb
WORKDIR /go/src/pgweb

RUN go get github.com/tools/godep
RUN godep restore
RUN godep go build && godep go install

CMD ["pgweb"]