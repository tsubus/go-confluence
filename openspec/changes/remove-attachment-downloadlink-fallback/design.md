## Context

Today attachment download has two public entry points with overlapping behavior:

- `Download(ctx, id, params, writer)` resolves metadata and downloads by attachment ID.
- `DownloadByURL(ctx, downloadURL, writer)` streams a caller-provided URL.

The current `Download` implementation also uses `DownloadByURL` internally as a fallback when the supported API endpoint fails before any bytes are written.

That creates an ambiguous contract:

```text
Download(id)
  ├─ supported API download by ID
  └─ fallback direct link download from metadata
```

After this change the contracts are separated:

```text
Download(id)
  └─ supported API download by ID only

DownloadByURL(url)
  └─ deprecated direct URL helper for caller-supplied URLs only
```

## Goals

- Make `Download(id)` use only the supported Confluence attachment download flow.
- Remove hidden fallback behavior that can mask failures in the supported path.
- Preserve a temporary escape hatch for callers that already have a URL, while clearly marking it deprecated.
- Warn at runtime when deprecated direct URL download is used.

## Non-Goals

- Removing `DownloadByURL` in this change.
- Changing method signatures.
- Reworking attachment metadata fetching or redirect-host validation beyond what is needed for the contract change.

## Design Decisions

### 1. `Download` stops after supported API download failure

If metadata lookup succeeds and a page/container ID is available, `Download` attempts the attachment API download endpoint.

- On success: return success.
- On failure before any bytes are written: return that API download failure directly.
- On failure after bytes are written: return that streaming failure directly, as today.
- Do not inspect or use `attachment.DownloadLink` as an alternate transport.

This keeps `Download` honest about the supported retrieval path.

### 2. Keep `DownloadByURL` temporarily, but deprecate it

`DownloadByURL` remains public for compatibility and for advanced callers that already possess a working URL. However:

- the Go doc comment should use the standard `Deprecated:` prefix so editors and `go doc` surface it
- package docs should reflect that it is deprecated and not used internally as fallback
- the method is treated as transitional API surface, not the recommended retrieval path

### 3. Runtime warning mechanism

The warning should be lightweight and non-fatal. The main design requirement is that it is visible enough to nudge callers away from the API without changing method signatures.

Preferred shape:

- emit one process-level warning per call site or per process invocation
- write to the standard library logger (`log.Printf`) with a concise message
- message should name `AttachmentService.DownloadByURL` and point users to `AttachmentService.Download` where applicable

Example message shape:

```text
confluence: AttachmentService.DownloadByURL is deprecated and will be removed in a future release; prefer AttachmentService.Download for Confluence attachment IDs
```

This is intentionally simple:
- no new client option
- no callback hook
- no error return changes

If implementation details suggest once-per-process suppression is easy, that is preferable to warning on every call.

## Risks

### Compatibility risk

Some users may unknowingly depend on the metadata-link fallback because their environments fail on the supported API endpoint but still allow direct link download. After this change, those calls will fail.

That is an intentional behavior correction, but it should be called out in proposal/docs so it is not surprising.

### Warning noise

If runtime warnings are emitted on every `DownloadByURL` call, noisy callers may see repeated log output. If straightforward, suppress repeated warnings after the first emission.

## Verification Strategy

- Update unit tests so `Download` no longer falls back to `downloadLink` when API-by-ID fails.
- Keep coverage for partial-stream failure behavior.
- Add focused coverage for the `DownloadByURL` deprecation marker/warning mechanism.
- Update public docs to reflect the new split between supported ID-based download and deprecated direct URL download.
