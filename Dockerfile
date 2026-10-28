FROM golang:1.27.3 AS builder

ENV MODULE=github.com/soerenschneider/acmevault
ENV CGO_ENABLED=0

WORKDIR /build/
COPY go.mod go.sum /build/
RUN go mod download

COPY . /build/
RUN go build -trimpath -buildvcs=false -ldflags="-w -X $MODULE/internal.BuildVersion=$(git describe --tags --abbrev=0 || echo dev) \
     -X $MODULE/internal.CommitHash=$(git rev-parse HEAD)" -o acmevault ./cmd

FROM gcr.io/distroless/static
COPY --from=builder /build/acmevault /acmevault
ENTRYPOINT ["/acmevault"]
