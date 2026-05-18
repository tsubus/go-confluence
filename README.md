# go-confluence

A Go client library for [Confluence](https://www.atlassian.com/software/confluence) REST v2 APIs and PDF export workflows.

- Supports **Confluence Server/Data Center** redirect-based PDF exports and **Confluence Cloud** background export tasks.
- Provides typed access to REST v2 endpoints for pages, spaces, attachments, and labels.
- Supports streaming downloads for PDFs and attachments to avoid unnecessary buffering.
- Ships with auto-generated API reference documentation (`DOCS.md`) in every package.

## Installation

```bash
go get github.com/umats/go-confluence
```

Requires Go 1.26 or later.

## Quick start

Create a client, then access the service you need:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    confluence "github.com/umats/go-confluence"
)

func main() {
    ctx := context.Background()

    client, err := confluence.NewClient(
        "https://wiki.example.com",
        confluence.WithBasicAuth("username", "password-or-api-token"),
        confluence.WithTimeout(60*time.Second),
    )
    if err != nil {
        log.Fatalf("create client: %v", err)
    }

    page, err := client.Page().Get(ctx, "123456789", nil)
    if err != nil {
        log.Fatalf("get page: %v", err)
    }

    fmt.Printf("Loaded page %q\n", page.Title)

    file, err := os.Create("page.pdf")
    if err != nil {
        log.Fatalf("create file: %v", err)
    }
    defer file.Close()

    if err := client.Export().PageTo(ctx, "123456789", file); err != nil {
        log.Fatalf("export page: %v", err)
    }

    fmt.Println("PDF downloaded successfully")
}
```

## Creating a client

[`NewClient`](client.go:148) validates the base URL, disables automatic redirect following so downloads can be checked explicitly, and returns service accessors for pages, spaces, attachments, labels, and export workflows.

```go
client, err := confluence.NewClient(
    "https://wiki.example.com",
    confluence.WithBasicAuth("username", "password-or-api-token"),
    confluence.WithRequireHTTPS(),
    confluence.WithAllowedRedirectHosts("downloads.example.com"),
)
if err != nil {
    log.Fatal(err)
}
```

### Client options

| Option | Description |
|--------|-------------|
| [`WithBasicAuth(username, password)`](client.go:57) | Sets username and password (or API token) for HTTP Basic Auth. |
| [`WithHTTPClient(hc)`](client.go:66) | Replaces the default `http.Client`. |
| [`WithPollInterval(d)`](client.go:77) | Sets how often to poll a Cloud export task (default: 3s). |
| [`WithPollTimeout(d)`](client.go:88) | Sets the maximum time to wait for a Cloud export task. |
| [`WithTimeout(d)`](client.go:99) | Sets the HTTP request timeout. |
| [`WithRequireHTTPS()`](client.go:110) | Enforces that `baseURL` uses the `https` scheme. |
| [`WithAllowedRedirectHosts(hosts...)`](client.go:118) | Allows downloads to follow redirects only to the base host plus these extra hosts. |
| [`WithAllowCrossHostContentURL()`](client.go:140) | Allows custom content URLs to point at a host different from the configured base URL. |

## Exporting pages as PDF

The [`ExportService`](export/service.go:1) handles PDF export for both Confluence Server/Data Center and Confluence Cloud:

- **Server/Data Center**: the export action returns a `302` redirect to the PDF.
- **Cloud**: the export starts a background task; the client polls the task progress endpoint and downloads the PDF when ready.

### Export to memory

```go
pdfBytes, err := client.Export().Page(ctx, "123456789")
if err != nil {
    log.Fatalf("export page: %v", err)
}
// pdfBytes is a []byte containing the PDF.
```

### Stream to a writer

If you want to stream the PDF directly to a file or another `io.Writer` without buffering it in memory, use [`ExportService.PageTo`](export/service.go:25):

```go
file, err := os.Create("page.pdf")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

err = client.Export().PageTo(ctx, "123456789", file)
if err != nil {
    log.Fatalf("export page: %v", err)
}
```

### Export error behavior

Export errors include the full URL the client attempted to access during request creation, task polling, redirect handling, or PDF download. Sentinel errors are preserved, so callers can still use `errors.Is` with [`ErrMissingLocation`](client.go:25), [`ErrTaskFailed`](client.go:29), [`ErrTaskResultEmpty`](client.go:31), and [`ErrTaskIDNotFound`](client.go:33).

## Downloading attachments

The [`AttachmentService`](attachment/service.go:1) supports both typed REST v2 attachment operations and streaming attachment downloads.

### Download by attachment ID

[`AttachmentService.Download`](attachment/service.go:230) fetches attachment metadata, resolves the attachment container, requests Confluence's supported attachment download endpoint, follows the returned media redirect, and streams the file to an `io.Writer`.

```go
file, err := os.Create("attachment.zip")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

err = client.Attachment().Download(ctx, "att-123456", nil, file)
if err != nil {
    log.Fatalf("download attachment: %v", err)
}
```

You can also pass query parameters, for example to request a specific version:

```go
version := 2
params := &confluence.AttachmentGetParams{Version: &version}
err = client.Attachment().Download(ctx, "att-123456", params, file)
```

Errors from attachment downloads also include the attempted URL so redirect validation, unexpected status responses, and streaming failures are easier to diagnose.

### Download by direct URL

[`AttachmentService.DownloadByURL`](attachment/service.go:285) is deprecated and is no longer used as a fallback by [`AttachmentService.Download`](attachment/service.go:230). Prefer `Download` for Confluence attachment IDs. Use `DownloadByURL` only when the caller already has an absolute attachment download URL.

```go
err = client.Attachment().DownloadByURL(ctx, "https://wiki.example.com/download/attachments/123/file.zip", file)
if err != nil {
    log.Fatalf("download attachment: %v", err)
}
```

## Other services

The root [`confluence`](client.go:1) package is the entry point for creating clients and accessing services. Domain APIs and shared types live in subpackages:

| Package | Description | Docs |
|---------|-------------|------|
| [`export`](export/service.go:1) | PDF export workflow. | [`export/DOCS.md`](export/DOCS.md) |
| [`page`](page/service.go:1) | REST v2 page operations. | [`page/DOCS.md`](page/DOCS.md) |
| [`space`](space/service.go:1) | REST v2 space operations. | [`space/DOCS.md`](space/DOCS.md) |
| [`attachment`](attachment/service.go:1) | REST v2 attachment operations and download helpers. | [`attachment/DOCS.md`](attachment/DOCS.md) |
| [`label`](label/service.go:1) | REST v2 label operations. | [`label/DOCS.md`](label/DOCS.md) |
| [`models`](models/models_generated.go:1) | OpenAPI-generated request/response models. | [`models/DOCS.md`](models/DOCS.md) |
| [`params`](params/params.go:1) | Query parameter structs and paginated result helpers. | [`params/DOCS.md`](params/DOCS.md) |
| [`internal/transport`](internal/transport/transport.go:1) | Internal transport, query, and JSON helpers. | [`internal/transport/DOCS.md`](internal/transport/DOCS.md) |

> **Note:** Every package includes a `DOCS.md` file generated by [gomarkdoc](https://github.com/princjef/gomarkdoc). These files contain the full API reference for types, functions, and methods in that package.

## Service interfaces

For dependency injection and testing, use the smaller interfaces in [`services_interfaces.go`](services_interfaces.go:1) such as [`PageReader`](services_interfaces.go:11), [`PageWriter`](services_interfaces.go:62), [`AttachmentReader`](services_interfaces.go:273), and [`ExportService`](services_interfaces.go:380). The larger [`PageService`](services_interfaces.go:172), [`AttachmentService`](services_interfaces.go:351), and related interfaces compose those subsets for backward compatibility.

## Query parameters

[`BuildQuery`](internal/transport/query.go:1) accepts structs that use `json` tags, `url.Values`, and `map[string]string`/`map[string][]string`. Struct fields respect `omitempty` and skip `json:"-"`.

## Development

This project uses [Task](https://taskfile.dev) for common commands:

| Command | Description |
|---------|-------------|
| `task test` | Run unit tests. |
| `task test:race` | Run unit tests with the race detector. |
| `task test:cover` | Run tests and print coverage. |
| `task test:integration` | Run integration tests with the `integration` build tag after loading `.env`. |
| `task lint` | Run `golangci-lint`. |
| `task lint:fix` | Run `golangci-lint` with auto-fix. |
| `task docs` | Regenerate package `DOCS.md` files with `gomarkdoc`. |

### Integration test environment

Integration tests fail fast when required credentials or fixture IDs are missing. Configure `.env` with:

```bash
CONFLUENCE_URL=https://wiki.example.com
CONFLUENCE_USERNAME=username
CONFLUENCE_PASSWORD=password-or-api-token
CONFLUENCE_PAGE_ID=123456789
CONFLUENCE_SPACE_ID=987654321
CONFLUENCE_ATTACHMENT_ID=attachment-id
CONFLUENCE_LABEL_ID=12345
# Optional: use a page that contains at least one attachment for attachment-focused scenarios.
# This may match CONFLUENCE_PAGE_ID only if that page has at least one attachment.
CONFLUENCE_ATTACHMENT_PAGE_ID=123456789
```

`CONFLUENCE_ATTACHMENT_PAGE_ID` is documented for integration fixture compatibility. If you use the same page as `CONFLUENCE_PAGE_ID`, that page must contain at least one attachment.