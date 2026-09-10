# Gravitee Automation SDKs

Go SDK clients and mock server for the Gravitee Access Management (AM) Automation API.

## Modules

| Module | Description |
|--------|-------------|
| `common` | Shared utilities: API context, response helpers, in-memory store, error types |
| `am` | Generated SDK clients for AM resources (domains, certificates, identity providers, reporters) |
| `am-mock-server` | Mock HTTP server implementing the AM Automation API, used for integration testing |
| `apim` | Placeholder for future APIM SDK |

## Prerequisites

- Go 1.26+

## Build & Test

```bash
# Run all tests
go test ./am/... ./am-mock-server/... ./common/...

# Run tests for a single module
go test ./am-mock-server/server/...

# Regenerate code (overlays + oapi-codegen)
go generate ./am/... ./am-mock-server/...
```

> **Note:** `go test ./...` does not work from the workspace root — specify module paths explicitly.

## Running the Mock Server

```bash
go run ./am-mock-server --port 8080
```

Options:

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | HTTP listen port |
| `--base-path` | `/automation` | API base path |
| `--auth-file` | | Path to auth config YAML (optional) |

## Code Generation

All generated code comes from the OpenAPI spec at `am/openapi/openapi.yaml`. Generation uses OpenAPI Overlay files to reshape the spec, then `oapi-codegen` to produce typed Go clients and strict servers.

Files ending in `.gen.go` are generated — do not edit them.

## License

[Apache License 2.0](LICENSE)
