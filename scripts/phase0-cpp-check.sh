#!/bin/sh
set -eu
# Native process-group verification is implemented as an opt-in Go integration test.
# macOS installations with mixed CLT/Xcode SDKs can set SDKROOT explicitly.
OIDESK_CPP_PROBE=1 go test -v -timeout 90s ./phase0 -run TestCPPProcessGroupTimeout
