MODULES := am apim common am-mock-server

HAS_GO = [ -n "$$(find $$mod -name '*.go' -print -quit)" ]

AM_OAS_BRANCH ?= master
AM_OAS_URL := https://raw.githubusercontent.com/gravitee-io/gravitee-access-management/refs/heads/$(AM_OAS_BRANCH)/docs/automation/openapi.yaml
AM_OAS_FILE := am/openapi/openapi.yaml

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
		-ignore "**/overlay.merged.yaml" \
		-ignore "am/pkg/sdk/overlay.yaml" \
		-ignore "am-mock-server/server/overlay.yaml" \
		-ignore ".github/**" \
		-ignore ".idea/**" \
		.

.PHONY: lint-fix
lint-fix: ## Auto-fix linting issues and add license headers
	@for mod in $(MODULES); do \
		find $$mod -name '*.go' -not -name '*.gen.go' -exec goimports -w {} +; \
	done
	@addlicense -f LICENSE_TEMPLATE.txt \
		-ignore "**/*.gen.go" \
		-ignore "**/overlay.merged.yaml" \
		-ignore "am/pkg/sdk/overlay.yaml" \
		-ignore "am-mock-server/server/overlay.yaml" \
		-ignore ".github/**" \
		-ignore ".idea/**" \
		.

##@ 🔄 Generate

.PHONY: sync-oas
sync-oas: ## Download AM Automation OAS from gravitee-access-management
	curl -fsSL "$(AM_OAS_URL)" -o $(AM_OAS_FILE)

.PHONY: check-oas
check-oas: ## Fail if committed OAS differs from upstream $(AM_OAS_BRANCH)
	curl -fsSL "$(AM_OAS_URL)" -o /tmp/am-openapi.yaml
	@diff -u $(AM_OAS_FILE) /tmp/am-openapi.yaml || { echo "OAS out of date. Run 'make sync-oas' and commit."; exit 1; }

.PHONY: generate
generate: ## Run go generate across all modules
	@for mod in $(MODULES); do \
		if $(HAS_GO); then \
			echo "==> generate $$mod ..."; \
			(cd $$mod && go generate ./...) || exit 1; \
		fi; \
	done

##@ 🧪 Test

.PHONY: test
test: ## Run tests
	@for mod in $(MODULES); do \
		if $(HAS_GO); then \
			echo "==> testing $$mod ..."; \
			if [ "$$mod" = "common" ]; then \
				(cd $$mod && go test -race ./...) || exit 1; \
			else \
				(cd $$mod && go test ./...) || exit 1; \
			fi; \
		fi; \
	done

##@ 🔧 Tools

.PHONY: tools
tools: ## Install development tools
	go install honnef.co/go/tools/cmd/staticcheck@v0.8.1
	go install github.com/mgechev/revive@v1.16.0
	go install github.com/google/addlicense@v1.2.0
	go install golang.org/x/tools/cmd/goimports@v0.50.0
