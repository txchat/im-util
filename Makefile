
help: ## Display this help screen
	@printf "Help doc:\nUsage: make [command]\n"
	@printf "[command]\n"
	@grep -h -E '^([a-zA-Z_-]|\%)+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: fmt_proto fmt_shell fmt_go

fmt: fmt_proto fmt_shell fmt_go

fmt_proto: ## format protobuf file
	@find . -name '*.proto' -not -path "./vendor/*" | xargs clang-format -i

fmt_shell: ## format shell file
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shfmt -w -s -i 4 -ci -bn

fmt_go: ## format go
	@find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -s -w
	@find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w

.PHONY: checkgofmt linter linter_test

checkgofmt: ## get all go files and run go fmt on them
	@files=$$(find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -l -s); if [ -n "$$files" ]; then \
		  echo "Error: 'make fmt' needs to be run on:"; \
		  find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -l -s ;\
		  exit 1; \
		  fi;
	@files=$$(find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w); if [ -n "$$files" ]; then \
		  echo "Error: 'make fmt' needs to be run on:"; \
		  find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w ;\
		  exit 1; \
		  fi;

#linter: ## Use gometalinter check code, ignore some unserious warning
#	@chmod +x ./script/golinter.sh
#	@./script/golinter.sh "filter"
#	@find . -name '*.sh' -not -path "./vendor/*" | xargs shellcheck

linter: ## Use gometalinter check code, ignore some unserious warning
	@golangci-lint run ./... && find . -name '*.sh' -not -path "./vendor/*" | xargs shellcheck

linter_test: ## Use gometalinter check code, for local test
	@chmod +x ./script/golinter.sh
	@./script/golinter.sh "test" "${p}"
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shellcheck
