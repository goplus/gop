ARG BASE_IMAGE=golang:1.22-bookworm

FROM $BASE_IMAGE AS build
ARG USE_GORELEASER_ARTIFACTS=0
ARG GORELEASER_ARTIFACTS_TARBALL
WORKDIR /usr/local/src/gop
COPY . .
ENV GOPROOT=/usr/local/gop
RUN set -eux; \
	mkdir $GOPROOT; \
	if [ $USE_GORELEASER_ARTIFACTS -eq 1 ]; then \
		tar -xzf "${GORELEASER_ARTIFACTS_TARBALL}" -C $GOPROOT; \
	else \
		git ls-tree --full-tree --name-only -r HEAD | grep -vE "^\." | xargs -I {} cp --parents {} $GOPROOT/; \
		./all.bash; \
		mv bin $GOPROOT/; \
	fi

FROM $BASE_IMAGE
ENV GOPROOT=/usr/local/gop
COPY --from=build $GOPROOT/ $GOPROOT/
ENV PATH=$GOPROOT/bin:$PATH
WORKDIR /gop
