all:
	cd cmd; go install -v ./...

release:
  # the release action expects the generated binaries placed in the path where make runs.
  # Refer: https://github.com/marketplace/actions/go-release-binaries
	go build -o gop github.com/goplus/gop/cmd/gop
	go build -o gopfmt github.com/goplus/gop/cmd/gopfmt