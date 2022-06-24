# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.18 AS build

WORKDIR /build
ADD . /build

RUN go mod download
RUN make build

# ------------------------------------------------------------------------------
# Release Stage
# ------------------------------------------------------------------------------
FROM alpine:3.16

RUN \
  apk update && \
  apk add --no-cache ca-certificates openssl postgresql && \
  update-ca-certificates && \
  rm -rf /var/cache/apk/*

COPY --from=build /build/pgweb /usr/bin/pgweb

EXPOSE 8081

CMD ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8081"]
