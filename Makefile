build: version-info
	go build -ldflags="-X 'acmevault/internal.BuildVersion=${VERSION}' -X 'acmevault/internal.CommitHash=${COMMIT_HASH}'" -o acmevault-server cmd/server/server.go
	go build -ldflags="-X 'acmevault/internal.BuildVersion=${VERSION}' -X 'acmevault/internal.CommitHash=${COMMIT_HASH}'" -o acmevault-client cmd/client/client.go

release: build
	sha256sum acmevault-* > checksums.sha256

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
	$(eval COMMIT_HASH := $(shell git rev-parse --short HEAD))

