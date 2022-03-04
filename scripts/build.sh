#!/bin/sh

set -eu

TARGET=${TARGET:-"target"}

GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"

TARGET="$TARGET/${GOOS}-${GOARCH}/docker-pps"
if [ "${GOOS}" = "windows" ]; then
    TARGET="${TARGET}.exe"
fi

case "$1" in
  fmt)
    go fmt $(go list)
  ;;
  clean)
    rm -rf ./target
  ;;
  build)
    go build -o "${TARGET}"
  ;;
  install)
    go install
  ;;
  *)
    echo "Usage: $0 {fmt|clean|build|install}"
    exit 1
esac

