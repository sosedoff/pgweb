# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.20-bookworm AS build

WORKDIR /build

RUN git config --global --add safe.directory /build
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile main.go ./
COPY static/ static/
COPY pkg/ pkg/
COPY .git/ .
RUN make build

# ------------------------------------------------------------------------------
# Release Stage
# ------------------------------------------------------------------------------
FROM gcr.io/distroless/base-debian12

COPY --from=build /build/pgweb /usr/bin/pgweb

EXPOSE 8081
ENTRYPOINT ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8081"]
