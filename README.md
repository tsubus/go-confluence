# go-confluence

A Go client library for [Confluence](https://www.atlassian.com/software/confluence) REST v2 APIs and PDF exports.

- Supports **Confluence Server/Data Center** and **Confluence Cloud** PDF exports.
- Provides typed access to REST v2 endpoints for pages, spaces, attachments, and labels.
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

    // Export a page by ID into memory.
    pdf, err := client.Export().Page(ctx, "123456789")
    if err != nil {
        log.Fatalf("export page: %v", err)
    }

    if err := os.WriteFile("page.pdf", pdf, 0o644); err != nil {
        log.Fatalf("write file: %v", err)
    }

    fmt.Println("PDF downloaded successfully")
}
```

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

If you want to stream the PDF directly to a file or another [`io.Writer`](export/service.go:25) without buffering it in memory, use [`ExportService.PageTo`](export/service.go:25):

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

### Export options

| Option | Description |
|--------|-------------|
| [`WithPollInterval(d)`](client.go:82) | How often to poll a Cloud export task (default: 3s). |
| [`WithPollTimeout(d)`](client.go:93) | Maximum time to wait for a Cloud export task. |
| [`WithAllowedRedirectHosts(hosts...)`](client.go:123) | Limits export download redirects to the allowed hosts. |

## Downloading attachments

The [`AttachmentService`](attachment/service.go:1) provides two ways to download attachment files:

### Download by attachment ID

[`AttachmentService.Download`](attachment/service.go:224) fetches the attachment metadata, resolves the download link, and streams the file to an [`io.Writer`](attachment/service.go:228):

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

You can also pass query parameters (for example, to request a specific version):

```go
params := &confluence.AttachmentGetParams{Version: 2}
err = client.Attachment().Download(ctx, "att-123456", params, file)
```

### Download by direct URL

If you already have a download URL (for example, from a page or space response), use [`AttachmentService.DownloadByURL`](attachment/service.go:242):

```go
err = client.Attachment().DownloadByURL(ctx, "https://wiki.example.com/download/attachments/123/file.zip", file)
if err != nil {
    log.Fatalf("download attachment: %v", err)
}
```

## Other services

The root [`confluence`](client.go:1) package is the entrypoint for creating clients and accessing services. Domain APIs and shared types live in subpackages:

| Package | Description | Docs |
|---------|-------------|------|
| [`export`](export/service.go:1) | PDF export workflow. | [`export/DOCS.md`](export/DOCS.md) |
| [`page`](page/service.go:1) | REST v2 page operations. | [`page/DOCS.md`](page/DOCS.md) |
| [`space`](space/service.go:1) | REST v2 space operations. | [`space/DOCS.md`](space/DOCS.md) |
| [`attachment`](attachment/service.go:1) | REST v2 attachment operations. | [`attachment/DOCS.md`](attachment/DOCS.md) |
| [`label`](label/service.go:1) | REST v2 label operations. | [`label/DOCS.md`](label/DOCS.md) |
| [`models`](models/models_generated.go:1) | OpenAPI-generated request/response models. | [`models/DOCS.md`](models/DOCS.md) |
| [`params`](params/params.go:1) | Query parameter structs and helpers. | [`params/DOCS.md`](params/DOCS.md) |
| [`internal/transport`](internal/transport/transport.go:1) | Internal transport helpers. | [`internal/transport/DOCS.md`](internal/transport/DOCS.md) |

> **Note:** Every package includes a `DOCS.md` file generated by [gomarkdoc](https://github.com/princjef/gomarkdoc). These files contain the full API reference for types, functions, and methods in that package.

### Optional subpackage usage

You can also import subpackages directly if you prefer a more modular setup:

```go
import (
    "context"
    "time"

    confluence "github.com/umats/go-confluence"
    "github.com/umats/go-confluence/export"
)

func main() {
    client, _ := confluence.NewClient("https://wiki.example.com", confluence.WithTimeout(30*time.Second))
    _ = export.NewService(client)
    _ = context.Background()
}
```

## Client options

| Option | Description |
|--------|-------------|
| [`WithBasicAuth(username, password)`](client.go:62) | Sets username and password (or API token) for HTTP Basic Auth. |
| [`WithHTTPClient(hc)`](client.go:71) | Replaces the default [`http.Client`](client.go:43). |
| [`WithPollInterval(d)`](client.go:82) | Sets how often to poll a Cloud export task (default: 3s). |
| [`WithPollTimeout(d)`](client.go:93) | Sets the maximum time to wait for a Cloud export task. |
| [`WithTimeout(d)`](client.go:104) | Sets the HTTP request timeout. |
| [`WithRequireHTTPS()`](client.go:115) | Enforces that `baseURL` uses the `https` scheme. |
| [`WithAllowedRedirectHosts(hosts...)`](client.go:123) | Limits export download redirects to the allowed hosts. |
| [`WithAllowCrossHostContentURL()`](client.go:144) | Allows custom content URLs to point at other hosts. |

## Service interfaces

For dependency injection and testing, use the smaller interfaces in [`services_interfaces.go`](services_interfaces.go:1) such as [`PageReader`](services_interfaces.go:224), [`PageWriter`](services_interfaces.go:275), and [`PageCommentsReader`](services_interfaces.go:328). The full [`PageService`](services_interfaces.go:385) composes these subsets for backward compatibility.

## Query parameters

[`BuildQuery`](internal/transport/query.go:1) accepts structs (using `json` tags), `url.Values`, and `map[string]string`/`map[string][]string`. Struct fields respect `omitempty` and skip `json:"-"`.

## Error handling

| Error | Reason |
|-------|--------|
| [`ErrMissingLocation`](client.go:33) | The export response did not include a `Location` header. |
| [`ErrTaskIDNotFound`](client.go:38) | The Cloud export response did not include a task ID. |
| [`ErrTaskFailed`](client.go:29) | The Confluence export task failed. |
| [`ErrTaskResultEmpty`](client.go:31) | The task finished without a result URL. |

## Development

This project uses [Task](https://taskfile.dev) for common commands:

| Command | Description |
|---------|-------------|
| `task test` | Run unit tests. |
| `task test:race` | Run unit tests with the race detector. |
| `task test:cover` | Run tests and print coverage. |
| `task test:integration` | Run integration tests (requires a `.env` file). |
| `task lint` | Run `golangci-lint`. |
| `task lint:fix` | Run `golangci-lint` with auto-fix. |

Integration tests expect a `.env` file with Confluence credentials:

```bash
CONFLUENCE_URL=https://wiki.example.com
CONFLUENCE_USERNAME=username
CONFLUENCE_PASSWORD=password-or-api-token
CONFLUENCE_PAGE_ID=123456789
CONFLUENCE_SPACE_ID=987654321
CONFLUENCE_ATTACHMENT_ID=attachment-id
```
