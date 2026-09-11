# AGENTS.md — Gravitee Automation SDKs

Universal AI agent instructions for the Gravitee Automation SDKs repository.
Consumed by: Claude Code, OpenAI Codex, GitHub Copilot, Cursor, Gemini/Jules, Windsurf, Zed, Aider, and others.

---

# Agent Guide

If you are an AI agent operating in this repository:
- You MUST read this file fully before making changes.
- You MUST follow the rules defined here.
- If any instruction conflicts with other files, this file takes precedence.

## 1. Quick Reference Commands

### Prerequisites

- Go 1.26+

### Build & Test

```bash
# Install lint/dev tools (staticcheck, revive, addlicense, goimports)
make tools

# Sync AM Automation OAS from gravitee-access-management (AM_OAS_BRANCH, default master)
make sync-oas

# Regenerate all code (overlays + oapi-codegen)
make generate

# Run all linters (vet, staticcheck, revive, license headers)
make lint

# Auto-fix lint issues and add license headers
make lint-fix

# Run all tests across all modules
make test

# Run tests for a single module
go test ./am-mock-server/server/...

# Run a single test
go test ./am-mock-server/server/... -run TestGetDomain404SDK

# Run the mock server
go run ./am-mock-server --port 8080
```

---

## 2. Project Context

### What Is This?

Go SDK clients and a mock server for the Gravitee Access Management (AM) Automation API. The SDK is generated from an OpenAPI spec using `oapi-codegen`, with OpenAPI Overlay files applied to reshape the spec before generation.

### Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.26 |
| Build | Go workspace (`go.work`) |
| Code generation | oapi-codegen + OpenAPI Overlay 1.1.0 |
| HTTP (mock server) | chi v5 |
| CLI (mock server) | cobra |
| Tests | stdlib `testing` + testify |

### Modules

A Go workspace with four modules:

| Module | Role |
|--------|------|
| `common` | Shared utilities: `apicontext` (auth + base URL), `response` (status helpers + generic `Payload[T]` extractor), `store` (channel-based in-memory generic store), `errors`, `refs` |
| `am-sdk` | Generated SDK clients for AM resources (domains, certificates, identity providers, reporters). Each resource lives in `am-sdk/pkg/sdk/<resource>/` with a `generate.go` and `cfg.yaml` |
| `am-mock-server` | Standalone mock HTTP server implementing the same OpenAPI spec with strict-server codegen. Used for integration-testing the SDK |
| `apim-sdk` | Placeholder module for a future APIM SDK (empty) |

---

## 3. Code Generation Pipeline

All generated code comes from a single OpenAPI spec at `am-sdk/openapi/openapi.yaml`. Generation is driven by `//go:generate` directives in `generate.go` files and involves two steps:

1. **Overlay merge** — `am-sdk/overlays/mergeoverlay.go` is a CLI tool that merges multiple YAML overlay files into one. Overlays rename models (strip `Automation` prefix via `x-go-type-name`), standardize operationIds (`get`, `list`, `upsert`, `delete`), and (for SDK clients) rewrite paths to bake `orgId`/`envId` into the server URL.
2. **oapi-codegen** — Generates typed Go clients or strict servers from the overlaid spec. Each package has a `cfg.yaml` controlling what gets generated and which tags to include.

### Overlay Chains

The overlay chain differs between SDK and mock server:

| Target | Overlays merged | What is generated |
|--------|----------------|-------------------|
| SDK (`am-sdk/pkg/sdk/<resource>/`) | `models.yaml` + `operations.yaml` + `overlay-paths.yaml` + per-resource `overlay.yaml` | Client + models, filtered by tag |
| Mock server (`am-mock-server/server/`) | `models.yaml` + `operations.yaml` + its own `overlay.yaml` | chi strict-server + models (all tags, original paths kept) |
| CRD models (`am-sdk/pkg/crd/domain/`) | Separate chain writing `overlay.merged.yaml` | Models only |

**Files ending in `.gen.go` are generated — do not edit them.**

---

## 4. Key Patterns

- **`store.Identifiable` interface** — Any type stored in `store.Store[T]` must implement `Identity() string`. The mock server's `identifiers.go` wires generated types to this interface.
- **`response.Payload[T]`** — Generic extractor that pulls `JSON200` from oapi-codegen response structs via reflection. Used in tests and intended for SDK consumers.
- **Test helpers** — `am-mock-server/server/helpers_test.go` contains reusable assertion functions (`assertSDKOK`, `assertGet404`, etc.) used across all resource test files. Each resource's tests exercise both raw HTTP and SDK client paths.

---

## 5. Conventions

- Search for similar implementations first before writing new code.
- Prefer existing patterns over inventing new ones.
- Keep changes small and reviewable; avoid unrelated refactors or reformatting.
