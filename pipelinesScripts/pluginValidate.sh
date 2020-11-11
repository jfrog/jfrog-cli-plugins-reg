#!/bin/bash
set -eu

cd pipelinesScripts/validator
go build
go run ./ "$@"; cd -
