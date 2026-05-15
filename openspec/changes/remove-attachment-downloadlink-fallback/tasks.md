## 1. Remove Attachment Download Fallback

- [x] 1.1 Update `attachment/service.go` so `AttachmentService.Download` returns the supported API download failure directly and no longer falls back to `attachment.DownloadLink`.
- [x] 1.2 Preserve the existing behavior that partial-stream failures from the supported API path are returned without attempting any alternate download path.
- [x] 1.3 Remove dead code that exists only to support metadata-link fallback if it is no longer used.

## 2. Deprecate Direct URL Download

- [x] 2.1 Mark `AttachmentService.DownloadByURL` as deprecated in Go doc comments so tooling surfaces the deprecation.
- [x] 2.2 Add a runtime deprecation warning for `DownloadByURL` with a concise migration message.
- [x] 2.3 Update generated/package docs to describe `DownloadByURL` as deprecated and no longer used internally as fallback.

## 3. Verify Behavior and Documentation

- [x] 3.1 Replace fallback-specific tests with assertions that `Download` no longer uses `downloadLink` when the supported API path fails.
- [x] 3.2 Keep or add coverage for partial-stream failure behavior on the supported API path.
- [x] 3.3 Add or update tests for the deprecation warning behavior of `DownloadByURL`.
- [x] 3.4 Update README examples/text so `Download` is presented as the supported attachment-ID flow and `DownloadByURL` is explicitly deprecated.
