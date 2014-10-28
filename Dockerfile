FROM golang:1.3

ENV PGHOST localhost
ENV PGPORT 5432
ENV PGUSER postgres
ENV PGDATABASE postgres

ADD . /go/src/pgweb
WORKDIR /go/src/pgweb

RUN touch Makefile
RUN make setup
RUN make dev

CMD pgweb --host $PGHOST --port $PGPORT --user $PGUSER --db $PGDATABASE
