#!/usr/bin/env bash

set -eo pipefail

echo "Current user: $(whoami)"
echo "Current directory: $(pwd)"
echo "Current PATH: $PATH"

# Check if we're in a container
if [ -f /.dockerenv ]; then
    echo "Running inside Docker container"
elif [ -f /run/.containerenv ]; then
    echo "Running inside container (Podman)"
else
    echo "Running on host system"
fi

# Force using Go 1.24
GO_BIN="/usr/local/go/bin/go"

echo "Checking Go binary:"
ls -l $GO_BIN
echo "Go version from direct path:"
/usr/local/go/bin/go version
echo "Go version from variable:"
$GO_BIN version

# Verify Go version matches go.mod
REQUIRED_VERSION="go1.24.0"
CURRENT_VERSION=$($GO_BIN version | awk '{print $3}')
if [[ "$CURRENT_VERSION" != "$REQUIRED_VERSION" ]]; then
    echo "Error: Go version mismatch. Required: $REQUIRED_VERSION, Found: $CURRENT_VERSION"
    echo "Please ensure you have Go 1.24.0 installed at /usr/local/go/bin/go"
    echo "Current PATH: $PATH"
    exit 1
fi

# Ensure we're using the right Go version for all commands
export PATH="/usr/local/go/bin:$PATH"

# Add GOPATH/bin to PATH
export PATH="$($GO_BIN env GOPATH)/bin:$PATH"

# Install required protobuf plugins
echo "Installing protobuf plugins..."
$GO_BIN install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
$GO_BIN install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest
$GO_BIN install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
$GO_BIN install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

protoc_install_proto_gen_doc() {
  echo "Installing protobuf protoc-gen-doc plugin"
  ($GO_BIN install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest 2> /dev/null)
}

echo "Generating gogo proto code"
cd proto

# Add the cosmos-sdk proto files to the include path
COSMOS_SDK_DIR=$($GO_BIN list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
export BUF_INCLUDE_PATH="$COSMOS_SDK_DIR/proto:$COSMOS_SDK_DIR/third_party/proto"

# Find all proto files in the current directory
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    echo $file
    if grep "option go_package" $file &> /dev/null ; then
      echo "buf generate --template buf.gen.gogo.yml " $file
      buf generate --template buf.gen.gogo.yml $file
    fi
  done
done

protoc_install_proto_gen_doc

echo "Generating proto docs"
buf generate --template buf.gen.doc.yml

cd ..

# move proto files to the right places
cp -r github.com/Victor118/liquidity/* ./
rm -rf github.com