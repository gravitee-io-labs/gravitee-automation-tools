MODULES := am apim common am-mock-server

HAS_GO = [ -n "$$(find $$mod -name '*.go' -print -quit)" ]

##@ 🧹 Lint

.PHONY: lint
lint: lint-vet lint-staticcheck lint-revive lint-licenses ## Run all linters

.PHONY: lint-vet
lint-vet: ## Run go vet
	@for mod in $(MODULES); do \
		if $(HAS_GO); then \
			echo "==> go vet $$mod ..."; \
			(cd $$mod && go vet ./...) || exit 1; \
		fi; \
	done

.PHONY: lint-staticcheck
lint-staticcheck: ## Run staticcheck
	@for mod in $(MODULES); do \
		if $(HAS_GO); then \
			echo "==> staticcheck $$mod ..."; \
			(cd $$mod && staticcheck ./...) || exit 1; \
		fi; \
	done

.PHONY: lint-revive
lint-revive: ## Run revive
	@for mod in $(MODULES); do \
		if $(HAS_GO); then \
			echo "==> revive $$mod ..."; \
			(cd $$mod && revive -config ../.revive.toml -formatter friendly -exclude "./**/*.gen.go" ./...) || exit 1; \
		fi; \
	done

.PHONY: lint-licenses
lint-licenses: ## Check license headers
	@echo "Checking license headers ..."
	@addlicense -check -f LICENSE_TEMPLATE.txt \
		-ignore "**/*.gen.go" \
		-ignore ".github/**" \
		.

.PHONY: lint-fix
lint-fix: ## Auto-fix linting issues and add license headers
	@for mod in $(MODULES); do \
		find $$mod -name '*.go' -not -name '*.gen.go' | xargs goimports -w; \
	done
	@addlicense -f LICENSE_TEMPLATE.txt \
		-ignore "**/*.gen.go" \
		-ignore ".github/**" \
		.

##@ 🧪 Test

.PHONY: test
test: ## Run tests
	@for mod in $(MODULES); do \
		if $(HAS_GO); then \
			echo "==> testing $$mod ..."; \
			(cd $$mod && go test ./...) || exit 1; \
		fi; \
	done

##@ 🔧 Tools

.PHONY: tools
tools: ## Install development tools
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/mgechev/revive@latest
	go install github.com/google/addlicense@latest
	go install golang.org/x/tools/cmd/goimports@latest
