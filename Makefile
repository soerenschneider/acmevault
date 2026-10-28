BUILD_DIR = builds
MODULE = github.com/soerenschneider/acmevault
BINARY_NAME = acmevault
CHECKSUM_FILE = $(BUILD_DIR)/checksum.sha256
SIGNATURE_KEYFILE = ~/.signify/github.sec
DOCKER_PREFIX = ghcr.io/soerenschneider

# Reproducible builds: always use the exact toolchain from go.mod, strip local paths and
# don't depend on VCS state (version and commit are injected via ldflags instead).
GO_VERSION := $(shell awk '/^go /{print $$2}' go.mod)
export GOTOOLCHAIN = go$(GO_VERSION)
GO_BUILD = CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-w -X '$(MODULE)/internal.BuildVersion=$(VERSION)' -X '$(MODULE)/internal.CommitHash=$(COMMIT_HASH)'"

tests:
	go test ./... -race -covermode=atomic -coverprofile=coverage.out
	go tool cover -html=coverage.out -o=coverage.html
	go tool cover -func=coverage.out -o=coverage.out

clean:
	git diff --quiet || { echo 'Dirty work tree' ; false; }
	rm -rf ./$(BUILD_DIR)

build: version-info
	$(GO_BUILD) -o $(BINARY_NAME) ./cmd

release: clean version-info cross-build
	cd $(BUILD_DIR) && shasum -a 256 $(BINARY_NAME)-* > $(notdir $(CHECKSUM_FILE))

signed-release: release
	pass keys/signify/github | signify -S -s $(SIGNATURE_KEYFILE) -m $(CHECKSUM_FILE)
	gh-upload-assets -o soerenschneider -r acmevault -f ~/.gh-token builds

cross-build: version-info
	GOOS=linux GOARCH=amd64       $(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64    ./cmd
	GOOS=linux GOARCH=arm GOARM=6 $(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-armv6    ./cmd
	GOOS=linux GOARCH=arm64       $(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-aarch64  ./cmd

# Builds all release binaries twice from separate copies of the source tree and compares checksums.
verify-reproducible: version-info
	$(eval TMP := $(shell mktemp -d))
	git worktree add --detach $(TMP)/a HEAD
	git worktree add --detach $(TMP)/b HEAD
	$(MAKE) -C $(TMP)/a cross-build VERSION=$(VERSION) COMMIT_HASH=$(COMMIT_HASH)
	$(MAKE) -C $(TMP)/b cross-build VERSION=$(VERSION) COMMIT_HASH=$(COMMIT_HASH)
	cd $(TMP)/a/$(BUILD_DIR) && shasum -a 256 * > ../../sums-a
	cd $(TMP)/b/$(BUILD_DIR) && shasum -a 256 * > ../../sums-b
	git worktree remove --force $(TMP)/a
	git worktree remove --force $(TMP)/b
	diff $(TMP)/sums-a $(TMP)/sums-b && echo "Build is reproducible" && cat $(TMP)/sums-a

docker-build:
	docker build -t "$(DOCKER_PREFIX)/acmevault-server" .

version-info:
	$(eval VERSION ?= $(shell git describe --tags --abbrev=0 || echo "dev"))
	$(eval COMMIT_HASH ?= $(shell git rev-parse HEAD))

fmt:
	find . -iname "*.go" -exec go fmt {} \; 

pre-commit-init:
	pre-commit install
	pre-commit install --hook-type commit-msg

pre-commit-update:
	pre-commit autoupdate

docs:
	rm -rf go-diagrams
	go run doc/main.go
	cd go-diagrams && dot -Tpng diagram.dot > ../overview.png
