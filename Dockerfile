ARG IMG_TAG=latest
ARG PLATFORM="linux/amd64"
ARG GO_VERSION="1.19"
ARG RUNNER_IMAGE="gcr.io/distroless/static"

FROM --platform=${PLATFORM} golang:${GO_VERSION}-alpine3.18 as builder
WORKDIR /src/app/
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

# From https://github.com/CosmWasm/wasmd/blob/master/Dockerfile
# For more details see https://github.com/CosmWasm/wasmvm#builds-of-libwasmvm
ARG ARCH=x86_64
# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.3.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.3.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9f5f387bddbd51a549f2d8eb2cc0f08030f311610c43004b96c141a733b681bc
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 129da0e50eaa3da093eb84c3a8f48da20b9573f3afdf83159b3984208c8ec5c3
RUN cp /lib/libwasmvm_muslc.${ARCH}.a /lib/libwasmvm_muslc.a

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES
RUN set -eux; apk add --no-cache ca-certificates build-base;

ARG VERSION=""
RUN BUILD_TAGS=muslc LINK_STATICALLY=true LDFLAGS=-buildid=$VERSION make build

# Add to a distroless container
ARG PLATFORM="linux/amd64"
FROM --platform=${PLATFORM} gcr.io/distroless/cc:$IMG_TAG
ARG IMG_TAG
COPY --from=builder /src/app/bin/chihuahuad /usr/local/bin/chihuahuad

EXPOSE 26656 26657 1317 9090

ENTRYPOINT ["chihuahuad"]
