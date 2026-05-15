## Why

`Attachment.Download` currently mixes two retrieval strategies:

1. fetch attachment metadata and try Confluence's attachment download endpoint by attachment ID
2. if that fails before any bytes are written, fall back to the attachment metadata `downloadLink`

That fallback is no longer a good contract. It hides the actual failure mode of the supported ID-based download path and couples the high-level API to a metadata link that is not a stable retrieval contract across recent Confluence behavior.

`DownloadByURL` also remains exposed as a public helper. Keeping it available for caller-supplied URLs can still be useful, but the library should stop relying on it internally for ID-based downloads and should clearly signal that direct URL download is deprecated.

## What Changes

- Make `AttachmentService.Download` API-endpoint-only after metadata lookup; remove fallback to metadata `downloadLink`.
- Preserve the current partial-stream safety rule: if the supported download path writes bytes and then fails, return that error directly.
- Mark `AttachmentService.DownloadByURL` as deprecated in Go doc comments and generated docs.
- Add a runtime deprecation warning when `DownloadByURL` is called so users are nudged away from the direct URL path.
- Update tests and documentation to reflect the separated contracts:
  - `Download(id, ...)` is the supported Confluence attachment retrieval path.
  - `DownloadByURL(url, ...)` is deprecated and not used internally as fallback.

## Capabilities

### Modified Capabilities
- `attachment-download`: Attachment download by ID uses only the supported attachment download endpoint and no longer falls back to metadata download links.
- `attachment-direct-url-download`: Direct URL attachment download remains available temporarily but is deprecated and emits a warning when used.

## Impact

- Affected code: `attachment/service.go`, root/package docs, and attachment download tests.
- Affected APIs:
  - `AttachmentService.Download` changes behavior by removing metadata-link fallback.
  - `AttachmentService.DownloadByURL` remains callable but is deprecated.
- User-visible behavior:
  - callers that previously succeeded only because of metadata-link fallback will now receive the supported API download failure.
  - callers using `DownloadByURL` will see a deprecation warning.
- Tests should assert that `Download` no longer falls back and that `DownloadByURL` deprecation is visible through the chosen warning mechanism.
