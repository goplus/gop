ARG BASE_IMAGE=golang:1.21-bookworm

FROM $BASE_IMAGE AS build
ARG USE_GORELEASER_ARTIFACTS=0
WORKDIR /usr/local/src/gop
COPY . .
ENV GOPROOT=/usr/local/gop
RUN set -eux; \
	mkdir -p $GOPROOT/bin; \
	git ls-tree --full-tree --name-only -r HEAD | grep -vE "^\." | xargs -I {} cp --parents {} $GOPROOT/; \
	if [ $USE_GORELEASER_ARTIFACTS -eq 1 ]; then \
		GOARCH=$(go env GOARCH); \
		BIN_DIR_SUFFIX=linux_$GOARCH; \
		[ $GOARCH = "amd64" ] && BIN_DIR_SUFFIX=${BIN_DIR_SUFFIX}_v1; \
		[ $GOARCH = "arm" ] && BIN_DIR_SUFFIX=${BIN_DIR_SUFFIX}_$(go env GOARM | cut -d , -f 1); \
		cp .dist/gop_$BIN_DIR_SUFFIX/bin/gop .dist/gopfmt_$BIN_DIR_SUFFIX/bin/gopfmt $GOPROOT/bin/; \
	else \
		./all.bash; \
		cp bin/gop bin/gopfmt $GOPROOT/bin/; \
	fi

FROM $BASE_IMAGE
ENV GOPROOT=/usr/local/gop
COPY --from=build $GOPROOT/ $GOPROOT/
ENV PATH=$GOPROOT/bin:$PATH
WORKDIR /gop
