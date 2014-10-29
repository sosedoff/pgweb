from golang:1.3

COPY . /go/src/github.com/sosedoff/pgweb
WORKDIR /go/src/github.com/sosedoff/pgweb
RUN make setup && make dev

ENTRYPOINT ["./pgweb"]

