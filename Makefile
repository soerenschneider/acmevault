BUILD_DIR = builds
MODULE = github.com/soerenschneider/acmevault
BINARY_NAME_SERVER = acmevault-server
BINARY_NAME_CLIENT = acmevault-client
CHECKSUM_FILE = $(BUILD_DIR)/checksum.sha256
SIGNATURE_KEYFILE = ~/.signify/github.sec
DOCKER_PREFIX = ghcr.io/soerenschneider

tests:
	go test ./...

clean:
	rm -rf ./$(BUILD_DIR)

build: version-info
	go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BINARY_NAME_SERVER) cmd/server/server.go
	go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BINARY_NAME_CLIENT) cmd/client/client.go

release: clean version-info cross-build-client cross-build-server
	sha256sum $(BUILD_DIR)/acmevault-* > $(CHECKSUM_FILE)

signed-release: release
	pass keys/signify/github | signify -S -s $(SIGNATURE_KEYFILE) -m $(CHECKSUM_FILE)

cross-build-server:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0       go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-linux-x86_64    cmd/server/server.go
	GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-linux-armv5     cmd/server/server.go
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-linux-armv6     cmd/server/server.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0       go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-linux-aarch64   cmd/server/server.go
	GOOS=openbsd GOARCH=amd64 CGO_ENABLED=0     go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-openbsd-x86_64  cmd/server/server.go

cross-build-client:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0       go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-linux-x86_64    cmd/client/client.go
	GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-linux-armv5     cmd/client/client.go
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-linux-armv6     cmd/client/client.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0       go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-linux-aarch64   cmd/client/client.go
	GOOS=openbsd GOARCH=amd64 CGO_ENABLED=0     go build -ldflags="-X '$(MODULE)/internal.BuildVersion=${VERSION}' -X '$(MODULE)/internal.CommitHash=${COMMIT_HASH}'" -o $(BUILD_DIR)/$(BINARY_NAME_CLIENT)-openbsd-x86_64  cmd/client/client.go

docker-build-server:
	docker build -t "$(DOCKER_PREFIX)/acmevault-server" --build-arg MODE=server .

docker-build-client:
	docker build -t "$(DOCKER_PREFIX)/acmevault-client" --build-arg MODE=client .

docker-build: docker-build-server docker-build-client

version-info:
	$(eval VERSION := $(shell git describe --tags --abbrev=0 || echo "dev"))
	$(eval COMMIT_HASH := $(shell git rev-parse HEAD))

fmt:
	find . -iname "*.go" -exec go fmt {} \; 

docs:
	rm -rf go-diagrams
	go run doc/main.go
	cd go-diagrams && dot -Tpng diagram.dot > ../overview.png
