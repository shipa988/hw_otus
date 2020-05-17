#!/usr/bin/env bash
set -xeuo pipefail

rm -f "$(command -v go-validate)"
rm -f ./models/*validation.go
rm -f ./models/*validation_generated.go

go install ./go-validate
go generate models/models.go
go test -v -tags generation ./models

echo "PASS"