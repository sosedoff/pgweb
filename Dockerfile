# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.18-buster AS build

WORKDIR /build
ADD . /build

RUN go mod download
RUN make build

# ------------------------------------------------------------------------------
# Release Stage
# ------------------------------------------------------------------------------
FROM debian:buster-slim

RUN \
  apt-get update && \
  apt-get install -y ca-certificates openssl postgresql netcat && \
  update-ca-certificates && \
  apt-get clean autoclean && \
  apt-get autoremove --yes && \
  rm -rf /var/lib/{apt,dpkg,cache,log}/

COPY --from=build /build/pgweb /usr/bin/pgweb

EXPOSE 8081

CMD ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8081"]
