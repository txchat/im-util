#!/bin/bash

# golangci-lint install
# https://golangci-lint.run/usage/install/
# darwin arm64需要从源码编译否则无法运行
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

# tool shellcheck install
# https://github.com/koalaman/shellcheck#installing
brew install shellcheck
