#!/usr/bin/env bash
# Install pre-commit and required dependencies.
set -e

echo "--Installing pre-commit & dependencies---"
brew install pre-commit
brew install golangci-lint
pre-commit install
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/go-critic/go-critic/cmd/gocritic@latest
go install golang.org/x/lint/golint@latest
go install github.com/BurntSushi/toml/cmd/tomlv@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/segmentio/golines@latest

touch  .git/hooks/commit-msg
echo "go-commitlinter" >> .git/hooks/commit-msg
chmod 755 .git/hooks/commit-msg
# shellcheck disable=SC2181
if [ "$?" = "0" ]; then
    echo "--- setup successful ---"
else
    echo "--- setup failed ---"
fi