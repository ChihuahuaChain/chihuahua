FROM golang:1.19-alpine AS go-builder

ARG E2E_SCRIPT_NAME
# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]
# we probably want to default to latest and error
# since this is predominantly for dev use
# hadolint ignore=DL3018
RUN set -eux; apk add --no-cache ca-certificates build-base;

# hadolint ignore=DL3018
RUN apk add git

# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers
# needed by github.com/zondax/hid
RUN apk add linux-headers
WORKDIR /chihuahua
COPY . /chihuahua

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9ecb037336bd56076573dc18c26631a9d2099a7f2b40dc04b6cae31ffb4c8f9a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 6e4de7ba9bad4ae9679c7f9ecf7e283dd0160e71567c6a7be6ae47c81ebe7f32

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.$(uname -m).a /lib/libwasmvm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
# then log output of file /code/bin/chihuahuad
# then ensure static linking
RUN E2E_SCRIPT_NAME=${E2E_SCRIPT_NAME} BUILD_TAGS=muslc LINK_STATICALLY=true make build-e2e-script 

# --------------------------------------------------------
FROM ubuntu

ARG E2E_SCRIPT_NAME

COPY --from=go-builder /chihuahua/build/${E2E_SCRIPT_NAME} /bin/${E2E_SCRIPT_NAME}
ENV HOME /chihuahua

WORKDIR $HOME

# rest server
EXPOSE 1317
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
#grpc
EXPOSE 9090
# Docker ARGs are not expanded in ENTRYPOINT in the exec mode. At the same time,
# it is impossible to add CMD arguments when running a container in the shell mode.
# As a workaround, we create the entrypoint.sh script to bypass these issues.
RUN echo "#!/bin/bash\n${E2E_SCRIPT_NAME} \"\$@\"" >> entrypoint.sh && chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]