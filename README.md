# go-confluence

A minimal Go client for exporting Confluence pages as PDFs. It works with both **Confluence Server/Data Center** and **Confluence Cloud**.

- **Server/Data Center**: the export endpoint returns a `302` redirect directly to the PDF.
- **Confluence Cloud**: the export starts a background task; the client polls the task progress endpoint and downloads the PDF when it is ready.

## Installation

```bash
go get github.com/umats/go-confluence
```

Requires Go 1.26 or later.

## Usage

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
    pdf, err := client.ExportPage(ctx, "123456789")
    if err != nil {
        log.Fatalf("export page: %v", err)
    }

    if err := os.WriteFile("page.pdf", pdf, 0o644); err != nil {
        log.Fatalf("write file: %v", err)
    }

    fmt.Println("PDF downloaded successfully")
}
```

### Stream to a writer

If you want to stream the PDF directly to a file or another `io.Writer` without buffering it in memory, use [`ExportPageTo`](export.go:28):

```go
file, err := os.Create("page.pdf")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

err = client.ExportPageTo(ctx, "123456789", file)
if err != nil {
    log.Fatalf("export page: %v", err)
}
```

## Client options

| Option | Description |
|--------|-------------|
| [`WithBasicAuth(username, password)`](client.go:60) | Sets username and password (or API token) for HTTP Basic Auth. |
| [`WithHTTPClient(hc)`](client.go:69) | Replaces the default [`http.Client`](client.go:43). |
| [`WithTimeout(d)`](client.go:102) | Sets the HTTP request timeout. |
| [`WithPollInterval(d)`](client.go:80) | Sets how often to poll a Cloud export task (default: 3s). |
| [`WithPollTimeout(d)`](client.go:91) | Sets the maximum time to wait for a Cloud export task. |
| [`WithRequireHTTPS()`](client.go:113) | Enforces that `baseURL` uses the `https` scheme. |

## Error handling

The client returns detailed errors for common failure scenarios:

- [`ErrMissingLocation`](client.go:31) — the export response lacked a `Location` header.
- [`ErrTaskFailed`](client.go:33) — the Confluence Cloud export task failed on the server.
- [`ErrTaskResultEmpty`](client.go:35) — the task finished but the result URL was empty.
- [`ErrTaskIDNotFound`](client.go:37) — the Cloud export HTML did not contain a task ID meta tag.

## Development

This project uses [Task](https://taskfile.dev) for common development tasks.

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
CONFLUENCE_BASE_URL=https://wiki.example.com
CONFLUENCE_USERNAME=username
CONFLUENCE_PASSWORD=password-or-api-token
CONFLUENCE_PAGE_ID=123456789
```

## License

MIT
