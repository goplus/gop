NAME := gop
RELEASE_VERSION := `git describe --tags`
BUILD_ROOT_DIR := build-dir

.PHONY: clean all
all: build

clean:
	rm -rf $(BUILD_ROOT_DIR)/*
	rm -f bin/*

build:
	go run cmd/make.go -build

dist:
	$(MAKE) clean
	mkdir -p bin/
	go build -o $(BUILD_ROOT_DIR)/make cmd/make.go
	$(MAKE) build-all

build-all: darwin-amd64.zip darwin-arm64.zip linux-386.zip linux-amd64.zip \
	linux-armv7.zip windows-386.zip windows-amd64.zip windows-armv7.zip windows-arm64.zip

build-dist:
	@mkdir -p bin/
	@rm -rf bin/*
	$(BUILD_ROOT_DIR)/make -build

%.zip: %
	@echo "Building $(NAME)-$(RELEASE_VERSION)-$@"

	@rm -f $(BUILD_ROOT_DIR)/$(NAME)-$(RELEASE_VERSION)-$@
	zip -r $(BUILD_ROOT_DIR)/$(NAME)-$(RELEASE_VERSION)-$@ . -x ".*" -x "*/.*" -x "$(BUILD_ROOT_DIR)/*"
	@echo "$(NAME)-$(RELEASE_VERSION)-$@ Done"

darwin-amd64:
	$(MAKE) GOARCH=amd64 GOOS=darwin BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

darwin-arm64:
	$(MAKE) GOARCH=arm64 GOOS=darwin BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

linux-386:
	$(MAKE) GOARCH=386 GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

linux-amd64:
	$(MAKE) GOARCH=amd64 GOOS=linux BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

linux-armv7:
	$(MAKE) GOARCH=arm GOOS=linux GOARM=7 BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

windows-386:
	$(MAKE) GOARCH=386 GOOS=windows EXE_SUFFIX=.exe BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

windows-amd64:
	$(MAKE) GOARCH=amd64 GOOS=windows EXE_SUFFIX=.exe BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

windows-armv7:
	$(MAKE) GOARCH=arm GOOS=windows EXE_SUFFIX=.exe GOARM=7 BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist

windows-arm64:
	$(MAKE) GOARCH=arm64 GOOS=windows EXE_SUFFIX=.exe BUILD_DIR=$(BUILD_ROOT_DIR)/$@/bin build-dist