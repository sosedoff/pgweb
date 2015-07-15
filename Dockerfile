FROM golang:1.4.2

COPY . /go/src/github.com/sosedoff/pgweb
WORKDIR /go/src/github.com/sosedoff/pgweb

RUN go get github.com/tools/godep

RUN godep restore
RUN godep go build && godep go install

EXPOSE 8081
CMD ["pgweb", "--bind", "0.0.0.0"]