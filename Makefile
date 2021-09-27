BUILD_DIR=builds

clean:
	rm -rf ./$(BUILD_DIR)

build: clean version-info
	go build -ldflags="-X 'acmevault/internal.BuildVersion=${VERSION}' -X 'acmevault/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/acmevault-server cmd/server/server.go
	go build -ldflags="-X 'acmevault/internal.BuildVersion=${VERSION}' -X 'acmevault/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/acmevault-client cmd/client/client.go

release: build
	sha256sum $(BUILD_DIR)/acmevault-* > $(BUILD_DIR)/checksums.sha256
	pass keys/signify/github | signify -S -s ~/.signify/github.sec -m $(BUILD_DIR)/checksums.sha256

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

