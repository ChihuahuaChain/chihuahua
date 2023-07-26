#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::') # grab everything after the space in "github.com/cometbft/cometbft v0.37.0"
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.7.0

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(chihuahua_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=chihuahua \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=chihuahuad \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION)

ifeq (cleveldb,$(findstring cleveldb,$(chihuahua_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (,$(findstring nostrip,$(chihuahua_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ifeq ($(LINK_STATICALLY),true)
        ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(chihuahua_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif
 
#$(info $$BUILD_FLAGS is [$(BUILD_FLAGS)])


all: check-go-version install

install: check-go-version go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/chihuahuad

build: check-go-version
	go build $(BUILD_FLAGS) -o bin/chihuahuad ./cmd/chihuahuad

BUILD_TARGETS := build install

build-reproducible-all: check-go-version build-reproducible-amd64 build-reproducible-arm64

build-reproducible-amd64: check-go-version
	ARCH=x86_64 PLATFORM=linux/amd64 $(MAKE) build-reproducible-generic

build-reproducible-arm64: check-go-version
	ARCH=aarch64 PLATFORM=linux/arm64 $(MAKE) build-reproducible-generic
	
build-reproducible-generic: check-go-version go.sum
	$(DOCKER) rm $(subst /,-,latest-build-$(PLATFORM)) || true
	DOCKER_BUILDKIT=1 $(DOCKER) build -t latest-build-$(PLATFORM) \
		--build-arg ARCH=$(ARCH) \
		--build-arg PLATFORM=$(PLATFORM) \
		--build-arg VERSION="$(VERSION)" \
		-f Dockerfile .
	$(DOCKER) create -ti --name $(subst /,-,latest-build-$(PLATFORM)) latest-build-$(PLATFORM) chihuahuad
	mkdir -p $(BUILDDIR)/$(NETWORK)/$(PLATFORM)/
	$(DOCKER) cp -a $(subst /,-,latest-build-$(PLATFORM)):/usr/local/bin/chihuahuad $(BUILDDIR)/$(NETWORK)/$(PLATFORM)/chihuahuad

# Add check to make sure we are using the proper Go version before proceeding with anything
check-go-version:
	@if ! go version | grep -q "go1.20"; then \
		echo "\033[0;31mERROR:\033[0m Go version 1.20 is required for compiling chihuahuad. It looks like you are using" "$(shell go version) \nThere are potential consensus-breaking changes that can occur when running binaries compiled with different versions of Go. Please download Go version 1.19 and retry. Thank you!"; \
		exit 1; \
	fi


###############################################################################
###                                Protobuf                                 ###
###############################################################################

containerProtoVer=0.13.0
containerProtoImage=ghcr.io/cosmos/proto-builder:$(containerProtoVer)

proto-gen:
	@echo "Generating Protobuf files"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; 
