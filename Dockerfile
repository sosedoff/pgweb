# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.22-bullseye AS build

# Set default build argument for CGO_ENABLED
ARG CGO_ENABLED=0
ENV CGO_ENABLED ${CGO_ENABLED}

WORKDIR /build

RUN git config --global --add safe.directory /build
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile main.go ./
COPY static/ static/
COPY pkg/ pkg/
RUN go build -ldflags="-X main.commit=render-build -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" 
# ------------------------------------------------------------------------------
# Fetch signing key
# ------------------------------------------------------------------------------
FROM debian:bullseye-slim AS keyring
ADD https://www.postgresql.org/media/keys/ACCC4CF8.asc keyring.asc
RUN apt-get update && \
    apt-get install -qq --no-install-recommends gpg
RUN gpg -o keyring.pgp --dearmor keyring.asc
RUN apt-get update && apt-get install -y git build-essential


# ------------------------------------------------------------------------------
# Release Stage
# ------------------------------------------------------------------------------
FROM debian:bullseye-slim

ARG keyring=/usr/share/keyrings/postgresql-archive-keyring.pgp
COPY --from=keyring /keyring.pgp $keyring
RUN . /etc/os-release && \
    echo "deb [signed-by=${keyring}] http://apt.postgresql.org/pub/repos/apt/ ${VERSION_CODENAME}-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    apt-get update && \
    apt-get install -qq --no-install-recommends ca-certificates openssl netcat curl postgresql-client

COPY --from=build /build/pgweb /usr/bin/pgweb

RUN useradd --uid 1000 --no-create-home --shell /bin/false pgweb
USER pgweb

EXPOSE 8080
ENTRYPOINT ["/pgweb", "--bind=0.0.0.0", "--listen=8080"]


# Set default port and binding


# Enable environment variable configuration
ENV DATABASE_URL="" PGWEB_READONLY=false
