# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.22-bullseye AS build

# Set default build argument for CGO_ENABLED
ARG CGO_ENABLED=0
ENV CGO_ENABLED ${CGO_ENABLED}

# Install essential build tools
RUN apt-get update && apt-get install -y git build-essential

WORKDIR /build

# Configure safe directory
RUN git config --global --add safe.directory /build

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application files
COPY Makefile main.go ./
COPY static/ static/
COPY pkg/ pkg/

# Build with timestamp-based versioning
RUN go build -ldflags="-X main.commit=render-build -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o /usr/bin/pgweb

# ------------------------------------------------------------------------------
# Fetch signing key
# ------------------------------------------------------------------------------
FROM debian:bullseye-slim AS keyring
ADD https://www.postgresql.org/media/keys/ACCC4CF8.asc keyring.asc
RUN apt-get update && \
    apt-get install -qq --no-install-recommends gpg && \
    gpg -o keyring.pgp --dearmor keyring.asc

# ------------------------------------------------------------------------------
# Release Stage
# ------------------------------------------------------------------------------
FROM debian:bullseye-slim

# Environment variables at top level
ENV DATABASE_URL="" PGWEB_READONLY=false

# Configure PostgreSQL repository
ARG keyring=/usr/share/keyrings/postgresql-archive-keyring.pgp
COPY --from=keyring /keyring.pgp $keyring
RUN . /etc/os-release && \
    echo "deb [signed-by=${keyring}] http://apt.postgresql.org/pub/repos/apt/ ${VERSION_CODENAME}-pgdg main" > /etc/apt/sources.list.d/pgdg.list

# Install dependencies
RUN apt-get update && \
    apt-get install -qq --no-install-recommends \
    ca-certificates \
    openssl \
    netcat \
    curl \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Copy pgweb binary
COPY --from=build /usr/bin/pgweb /usr/bin/pgweb

# Set permissions
RUN chmod +x /usr/bin/pgweb

# Create non-root user
RUN useradd --uid 1000 --no-create-home --shell /bin/false pgweb
USER pgweb

# Expose port and set entrypoint
EXPOSE 8080
ENTRYPOINT ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8080"]
