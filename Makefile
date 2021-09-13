build: version-info
	env CGO_ENABLED=0 go build -ldflags="-X 'acmevault/internal.BuildTime=${BUILD_TIME}' -X 'acmevault/internal.BuildVersion=${VERSION}' -X 'acmevault/internal.CommitHash=${COMMIT_HASH}'" -o acmevault-server cmd/server/server.go
	env CGO_ENABLED=0 go build -ldflags="-X 'acmevault/internal.BuildTime=${BUILD_TIME}' -X 'acmevault/internal.BuildVersion=${VERSION}' -X 'acmevault/internal.CommitHash=${COMMIT_HASH}'" -o acmevault-client cmd/client/client.go

fmt:
	find . -iname "*.go" -exec go fmt {} \; 

docs:
	rm -rf go-diagrams
	go run doc/main.go
	cd go-diagrams && dot -Tpng diagram.dot > ../overview.png

tests:
	go test ./...

version-info:
	$(eval VERSION := $(shell git describe --tags || echo "dev"))
	$(eval BUILD_TIME := $(shell date --rfc-3339=seconds))
	$(eval COMMIT_HASH := $(shell git rev-parse --short HEAD))

