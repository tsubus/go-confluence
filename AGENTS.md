# Repository Guidelines

## Project Overview

`github.com/umats/go-confluence` is a Go client library for Confluence REST v2 APIs and PDF export workflows. It provides typed services for pages, spaces, attachments, labels, and PDF export, supporting both Confluence Server/Data Center redirect-based exports and Confluence Cloud background export tasks.

## Architecture & Data Flow

- `confluence.NewClient(baseURL, opts...)` in `client.go` is the public entry point. It validates and normalizes the base URL, configures auth/timeouts/redirect policy, and stores shared HTTP configuration.
- Root service accessors in `client_services.go` (`Page()`, `Space()`, `Attachment()`, `Label()`, `Export()`) wrap the shared client into package-specific services.
- REST v2 services (`page`, `space`, `attachment`, `label`) use `internal/transport.V2Client.DoJSON` for URL construction, query encoding, JSON request bodies, response decoding, and non-2xx `APIError` handling.
- `params` contains typed query-parameter structs and generic `MultiEntityResult[T]` list wrappers. `internal/transport.BuildQuery` converts structs, maps, or `url.Values` into query strings using `json` tags.
- `models/models_generated.go` contains generated OpenAPI request/response types from `spec/openapi_subset.json`; do not hand-edit generated models.
- Export flow:
  - Server/Data Center: export endpoint returns `302`; client validates allowed redirect hosts and streams/downloads the PDF.
  - Cloud: export endpoint returns HTML with a task ID; client polls `/api/v2/pdfexporttask/progress/{taskID}` until a result URL is ready, then downloads it.
- Attachment download flow: `AttachmentService.Download` fetches metadata, asks Confluence for the supported download endpoint, follows the media redirect, validates hosts, and streams to an `io.Writer`.

## Key Directories

- `.`: root `confluence` package, client construction, service interfaces, compatibility aliases, top-level docs/tests.
- `internal/transport/`: internal HTTP, URL, query, JSON, and API error helpers shared by services.
- `page/`, `space/`, `attachment/`, `label/`: typed REST v2 service implementations.
- `export/`: PDF export implementation, redirect downloads, Cloud task polling, test-exposed helpers.
- `params/`: strongly typed query parameter structs and generic result wrappers.
- `models/`: generated model types; treat as generated code.
- `spec/`: OpenAPI subset used for generated models.
- `cmd/discover_ids/`: ignored helper command for discovering live Confluence fixture IDs from `CONFLUENCE_*` env vars.
- `openspec/`: spec/change workflow artifacts; useful for understanding accepted behavior changes.

## Development Commands

Prefer Task v3 targets for CI parity:

```bash
task test              # go test -v ./...
task test:race         # CGO_ENABLED=1 go test -race -v ./...
task test:cover        # go test -cover ./...
task test:integration  # loads .env; go test -tags=integration -v ./...
task lint              # golangci-lint run -v
task lint:fix          # golangci-lint run --fix -v
task docs              # regenerate package DOCS.md via gomarkdoc
```

Useful targeted examples:

```bash
go test . -run '^TestExportPage$' -v
go test ./... -run 'PageService_Methods'
go test ./internal/transport -run TestBuildQuery
go test -fuzz=FuzzExtractTaskID ./...
```

## Code Conventions & Common Patterns

- Go version is `1.26` (`go.mod`). Use standard Go formatting; lint is strict through `.golangci.yml` and CI uses `golangci-lint` v2.12.1.
- Public network methods take `context.Context` first and return typed `models`/`params` values plus `error`.
- Client configuration uses functional options: `WithBasicAuth`, `WithHTTPClient`, `WithTimeout`, `WithPollInterval`, `WithPollTimeout`, `WithRequireHTTPS`, `WithAllowedRedirectHosts`, `WithAllowCrossHostContentURL`.
- Services follow `NewService(client *transport.Client) *Service`; REST services hold a `*transport.V2Client`, while export uses `*transport.Client` directly.
- Prefer streaming APIs for large content: `Export.PageTo(ctx, pageID, io.Writer)` and `Attachment.Download(..., io.Writer)` avoid buffering; `Export.Page` intentionally returns `[]byte` for in-memory use.
- Error handling is contract-heavy: preserve sentinel errors (`ErrMissingLocation`, `ErrTaskIDNotFound`, `ErrTaskFailed`, `ErrTaskResultEmpty`) and include attempted URLs/status/body context where existing tests expect it.
- Security-sensitive redirect behavior lives in client/export/attachment transport code. The HTTP client deliberately avoids automatic redirects so code can inspect `Location` and validate allowed hosts.
- Dependency injection is via small interfaces in `services_interfaces.go` (`PageReader`, `PageWriter`, etc.) and `HTTPDoer` in `client.go`. Return concrete implementations from constructors; consume interfaces at boundaries.
- Backward compatibility aliases in `services_compat.go` expose `params` and `models` types from the root package. Update aliases when moving or adding public API types that consumers may need.
- Test-only hooks are explicit `ForTest` functions or exported helpers (`BuildQueryForTest`, `DoJSONForTest`, `DecodeResponseForTest`, `ExtractTaskIDForTest`). Avoid broadening production visibility unnecessarily.

## Important Files

- `client.go`: client options, validation, defaults, exported sentinel errors.
- `client_services.go`: root service accessors and interface assertions.
- `services_interfaces.go`: DI-friendly service interfaces.
- `services_compat.go`: root-package compatibility type aliases.
- `internal/transport/transport.go`: base request and URL construction.
- `internal/transport/v2_client.go`: REST v2 JSON request/response flow and `APIError`.
- `internal/transport/query.go`: reflection-based query encoding.
- `export/exporter.go`, `export/task.go`, `export/download.go`: PDF export workflows.
- `attachment/service.go`: attachment REST operations and streaming download behavior.
- `models/models_generated.go`: generated models; do not edit manually.
- `Taskfile.yml`: canonical local commands.
- `.github/workflows/ci.yml`: CI test, race, integration, and lint behavior.
- `.golangci.yml`: strict lint/format configuration.
- `README.md` and package `DOCS.md`: usage examples and generated API reference.

## Runtime/Tooling Preferences

- Runtime/toolchain: Go `1.26` or newer as declared in `go.mod`.
- Task runner: Task v3.x; use `task ...` targets rather than inventing ad hoc scripts.
- Package management: Go modules; no vendoring by default (`GO_VENDOR` in `Taskfile.yml` is intentionally blank).
- Docs tooling: `gomarkdoc` for package `DOCS.md` generation.
- Lint tooling: `golangci-lint` v2.12.1 for CI parity.
- Integration credentials should come from `.env` locally or GitHub secrets/vars in CI; never hard-code Confluence credentials or fixture IDs.

## Testing & QA

- Unit tests use standard `testing`, `httptest`, and `github.com/stretchr/testify/require`/`assert`.
- Test style is table-driven with `t.Run`, local `httptest.Server`/`http.ServeMux` fixtures, captured request structs, custom `RoundTripper`s, and deterministic failing readers/writers.
- Assert behavior, not just plumbing: request method/path/query/auth/body, streaming results, redirect validation, context cancellation, timeout behavior, sentinel errors, and URL/status context in failures.
- Integration tests are in `client_integration_test.go` behind `//go:build integration`; run with `task test:integration` or `go test -tags=integration ./...`.
- Required integration env vars include `CONFLUENCE_URL`, `CONFLUENCE_USERNAME`, `CONFLUENCE_PASSWORD`, `CONFLUENCE_PAGE_ID`, `CONFLUENCE_SPACE_ID`, `CONFLUENCE_ATTACHMENT_ID`, and `CONFLUENCE_LABEL_ID`; CI also provides attachment page fixture IDs.
- Existing tests intentionally fail fast on missing integration fixture data. Do not convert required live-fixture failures into skips unless the spec changes.
- Fuzz coverage exists for export task ID extraction: `go test -fuzz=FuzzExtractTaskID ./...`.
