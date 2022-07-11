#!/bin/bash

# golangci-lint install
# https://golangci-lint.run/usage/install/
brew install golangci-lint
brew upgrade golangci-lint

# shfmt install
# https://github.com/mvdan/sh
go install mvdan.cc/sh/v3/cmd/shfmt@latest

# clang-format install
# https://formulae.brew.sh/formula/clang-format
brew install clang-format

# goimports install
# https://pkg.go.dev/golang.org/x/tools/cmd/goimports
go install golang.org/x/tools/cmd/goimports@latest
