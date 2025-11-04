FROM golang:1.25 AS build
WORKDIR /build

ARG SHA
ARG TAG

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w \
        -X 'github.com/0xSplits/pulsar/pkg/runtime.sha=${SHA}' \
        -X 'github.com/0xSplits/pulsar/pkg/runtime.tag=${TAG}'" \
    -a \
    -o pulsar main.go



FROM gcr.io/distroless/static:nonroot
WORKDIR /image

COPY .env .env
COPY --from=build /build/pulsar .
USER 65532:65532

ENV PULSAR_HTTP_HOST="0.0.0.0"
ENV PULSAR_LOG_LEVEL="info"

ENTRYPOINT ["/image/pulsar"]
