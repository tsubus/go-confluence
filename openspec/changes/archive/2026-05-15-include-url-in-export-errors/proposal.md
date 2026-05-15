## Why

Errors raised while exporting pages or downloading attachments currently omit the full URL being requested in several failure paths. Including that URL makes failures from redirects, Cloud export polling, and direct downloads actionable without requiring request tracing or reproduction.

## What Changes

- Include the full attempted URL in errors from page export request creation/execution, unexpected export responses, and export response body handling.
- Include the full attempted URL in errors from PDF download request creation/execution, non-200 download responses, response streaming, and close/read failures.
- Include the full attempted URL in errors from Cloud export task polling, including request creation/execution, non-200 responses, response reads, JSON decoding, and polling timeout/cancellation context.
- Include the full attempted URL in attachment direct-download and redirect-download errors, including host validation, request creation/execution, non-200 responses, and response streaming.
- Preserve existing sentinel error behavior where callers use `errors.Is` for export conditions.
- No breaking API changes.

## Capabilities

### New Capabilities
- `export-error-url-context`: Error reporting for page export and attachment download operations includes the full URL being accessed when a request or download fails.

### Modified Capabilities

## Impact

- Affected code: `export` package page export flow, PDF download helper, Cloud export task polling, and `attachment` package download flow.
- Affected APIs: public errors returned by `ExportService.Page`, `ExportService.PageTo`, `AttachmentService.Download`, and `AttachmentService.DownloadByURL` gain additional context in their messages; method signatures stay unchanged.
- Tests should assert URL context is present on representative export, poll, and attachment download failures while preserving `errors.Is` behavior for sentinel errors.
