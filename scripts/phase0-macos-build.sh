#!/bin/sh
set -eu
if [ "$(uname -s)" != Darwin ]; then echo "macOS only" >&2; exit 1; fi
developer_dir="$(xcode-select -p)"
sdk="$developer_dir/Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk"
if [ ! -d "$sdk" ]; then echo "Xcode SDK missing: $sdk" >&2; exit 1; fi
export SDKROOT="$sdk"
wails build -clean -nocolour
