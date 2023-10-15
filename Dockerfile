FROM golang:1.21.2 as builder

ENV MODULE=github.com/soerenschneider/acmevault
ENV CGO_ENABLED=0

WORKDIR /build/
ADD go.mod go.sum /build
RUN go mod download

ADD . /build/
RUN go build -ldflags="-X $MODULE/internal.BuildVersion=$(git describe --tags --abbrev=0 || echo dev) \
     -X $MODULE/internal.CommitHash=$(git rev-parse HEAD)" -o acmevault "cmd/main.go"

FROM gcr.io/distroless/base
COPY --from=builder /build/acmevault /acmevault
ENTRYPOINT ["/acmevault"]
