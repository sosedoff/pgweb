# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.19-bullseye AS build

WORKDIR /build
ADD . /build

RUN git config --global --add safe.directory /build
RUN go mod download
RUN make build

# ------------------------------------------------------------------------------
# Release Stage
# ------------------------------------------------------------------------------
FROM debian:bullseye-slim

RUN \
  apt-get update && \
  apt-get install -y ca-certificates openssl netcat curl gnupg lsb-release && \
  update-ca-certificates

RUN \
  curl --silent https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add && \
  echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" | tee  /etc/apt/sources.list.d/pgdg.list && \
  apt-get update && apt-get install -y postgresql-client

RUN \
  apt-get clean autoclean && \
  apt-get autoremove --yes && \
  rm -rf /var/lib/{apt,dpkg,cache,log}/

COPY --from=build /build/pgweb /usr/bin/pgweb

EXPOSE 8081
ENTRYPOINT ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8081"]
