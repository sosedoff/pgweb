FROM golang:1.3

COPY . /go/src/pgweb
WORKDIR /go/src/pgweb

RUN touch Makefile
RUN make setup
RUN make dev

ENTRYPOINT ["pgweb"]
