all:
	cd cmd; go install -v ./...

release:
  # the release action expects the generated binaries placed in the path where make runs.
  # Refer: https://github.com/marketplace/actions/go-release-binaries
	CGO_ENABLED=0 go build -o gop -ldflags "-X 'github.com/goplus/gop.buildVersion=${RELEASE_VERSION}'" github.com/goplus/gop/cmd/gop
	CGO_ENABLED=0 go build -o gopfmt -ldflags "-X 'github.com/goplus/gop.buildVersion=${RELEASE_VERSION}'" github.com/goplus/gop/cmd/gopfmt