NAME := gop
PACKAGE_NAME := github.com/goplus/gop
RELEASE_VERSION := `git describe --tags`
COMMIT := `git rev-parse --short HEAD`
BUILD_DATE := `date -u +"%Y-%m-%d_%H-%M-%S"`

PLATFORM := linux
BUILD_DIR=.
BUILD_ROOT_DIR := build-dir
BUILD_LDFLAGS := -X 'github.com/goplus/gop/env.buildDate=$(BUILD_DATE)'
BUILD_LDFLAGS +=  -X '$(PACKAGE_NAME).buildCommit=$(COMMIT)'
BUILD_LDFLAGS +=  -X '$(PACKAGE_NAME).buildVersion=$(RELEASE_VERSION)'

EXE_SUFFIX=

.PHONY: clean all
all:
	cd cmd; go install -v ./...

clean:
	rm -rf $(BUILD_ROOT_DIR)
	rm -f *.zip
	rm -f gop gopfmt


release:
  # the release action expects the generated binaries placed in the path where make runs.
  # Refer: https://github.com/marketplace/actions/go-release-binaries
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/gop$(EXE_SUFFIX) -ldflags "$(BUILD_LDFLAGS)" github.com/goplus/gop/cmd/gop
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/gopfmt$(EXE_SUFFIX) -ldflags "$(BUILD_LDFLAGS)" github.com/goplus/gop/cmd/gopfmt

gop:
	mkdir -p $(BUILD_DIR)
	$(MAKE) release

release-all: darwin-amd64.zip darwin-arm64.zip linux-386.zip linux-amd64.zip \
	linux-arm.zip linux-armv5.zip linux-armv6.zip linux-armv7.zip linux-armv8.zip \
	linux-mips-softfloat.zip linux-mips-hardfloat.zip linux-mipsle-softfloat.zip linux-mipsle-hardfloat.zip \
	linux-mips64.zip linux-mips64le.zip freebsd-386.zip freebsd-amd64.zip \
	windows-386.zip windows-amd64.zip windows-arm.zip windows-armv6.zip windows-armv7.zip windows-arm64.zip

%.zip: %
	@echo "Building $(NAME)-$(RELEASE_VERSION)-$@"
	@mkdir -p $(BUILD_ROOT_DIR)/$</bin
	
	@rm -f $(NAME)-$(RELEASE_VERSION)-$@
	@cd $(BUILD_ROOT_DIR)/$</ && zip -r ../../$(NAME)-$(RELEASE_VERSION)-$@ bin
	@zip -dur $(NAME)-$(RELEASE_VERSION)-$@ builtin go.mod go.sum LICENSE README*.md
	@echo "$(NAME)-$(RELEASE_VERSION)-$@ Done"
	

darwin-amd64:
	$(MAKE) GOARCH=amd64 GOOS=darwin BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

darwin-arm64:
	$(MAKE) GOARCH=arm64 GOOS=darwin BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-386:
	$(MAKE) GOARCH=386 GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-amd64:
	$(MAKE) GOARCH=amd64 GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-arm:
	$(MAKE) GOARCH=arm GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-armv5:
	$(MAKE) GOARCH=arm GOOS=linux GOARM=5 BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-armv6:
	$(MAKE) GOARCH=arm GOOS=linux GOARM=6 BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-armv7:
	$(MAKE) GOARCH=arm GOOS=linux GOARM=7 BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-armv8:
	$(MAKE) GOARCH=arm64 GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-mips-softfloat:
	$(MAKE) GOARCH=mips GOMIPS=softfloat GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-mips-hardfloat:
	$(MAKE) GOARCH=mips GOMIPS=hardfloat GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-mipsle-softfloat:
	$(MAKE) GOARCH=mipsle GOMIPS=softfloat GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-mipsle-hardfloat:
	$(MAKE) GOARCH=mipsle GOMIPS=hardfloat GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-mips64:
	$(MAKE) GOARCH=mips64 GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

linux-mips64le:
	$(MAKE) GOARCH=mips64le GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

freebsd-386:
	$(MAKE) GOARCH=386 GOOS=freebsd BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

freebsd-amd64:
	$(MAKE) GOARCH=amd64 GOOS=freebsd BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

windows-386:
	$(MAKE) GOARCH=386 GOOS=windows EXE_SUFFIX=.exe BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

windows-amd64:
	$(MAKE) GOARCH=amd64 GOOS=windows EXE_SUFFIX=.exe BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

windows-arm:
	$(MAKE) GOARCH=arm GOOS=windows EXE_SUFFIX=.exe BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

windows-armv6:
	$(MAKE) GOARCH=arm GOOS=windows EXE_SUFFIX=.exe GOARM=6 BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

windows-armv7:
	$(MAKE) GOARCH=arm GOOS=windows EXE_SUFFIX=.exe GOARM=7 BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release

windows-arm64:
	$(MAKE) GOARCH=arm64 GOOS=windows EXE_SUFFIX=.exe BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin release