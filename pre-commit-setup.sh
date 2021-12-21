#!/usr/bin/env bash
# Install pre-commit and required dependencies.
set -e

echo "--Installing pre-commit & dependencies---"
brew install pre-commit
pre-commit install
GO111MODULE=off go get github.com/fzipp/gocyclo
GO111MODULE=off go get golang.org/x/tools/cmd/goimports
GO111MODULE=off go get -v -u github.com/go-critic/go-critic/cmd/gocritic
GO111MODULE=off go get -u golang.org/x/lint/golint

# shellcheck disable=SC2181
if [ "$?" = "0" ]; then
    echo "--- setup successful ---"
else
    echo "--- setup failed ---"
fi