# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.26-trixie AS build

# Set default build argument for CGO_ENABLED
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}

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
FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=build /build/pgweb /usr/bin/pgweb

USER nonroot

EXPOSE 8081

ENTRYPOINT ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8081"]
