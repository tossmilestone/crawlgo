#!/bin/sh

set -o errexit
set -o nounset

if [ -z "${PKG}" ]; then
    echo "PKG must be set"
    exit 1
fi
if [ -z "${ARCH}" ]; then
    echo "ARCH must be set"
    exit 1
fi
if [ -z "${VERSION}" ]; then
    echo "VERSION must be set"
    exit 1
fi

export CGO_ENABLED=0
export GOARCH="${ARCH}"

GO_SUBCMD="build"

if [ "${INSTALL}" = "y" ]; then
    GO_SUBCMD="install -installsuffix 'static'"
fi

go ${GO_SUBCMD}                                            \
    -ldflags "-X ${PKG}/pkg/version.VERSION=${VERSION}"    \
    ./...
