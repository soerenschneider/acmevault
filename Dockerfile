FROM golang:1.23.0 as builder

ENV MODULE=github.com/soerenschneider/acmevault
ENV CGO_ENABLED=0

WORKDIR /build/
ADD go.mod go.sum /build
RUN go mod download

ADD . /build/
RUN go build -ldflags="-w -X $MODULE/internal.BuildVersion=$(git describe --tags --abbrev=0 || echo dev) \
     -X $MODULE/internal.CommitHash=$(git rev-parse HEAD)" -o acmevault ./cmd

FROM gcr.io/distroless/base
COPY --from=builder /build/acmevault /acmevault
ENTRYPOINT ["/acmevault"]
