## 1. Export URL Error Context

- [x] 1.1 Add full export URL context to page export request creation, execution, response read, unexpected status, and response close errors in `export/exporter.go`.
- [x] 1.2 Add resolved redirect download URL context to missing `Location`, redirect parsing, host validation, request creation, execution, non-200 response, response read, stream, and close errors in `export/download.go`.
- [x] 1.3 Add full Cloud task progress URL context to poll request creation, execution, response read, non-200 status, JSON decode, and polling wait errors in `export/task.go`.
- [x] 1.4 Preserve `errors.Is` compatibility for existing export sentinel errors while adding URL context.

## 2. Attachment URL Error Context

- [x] 2.1 Add direct attachment download URL context to URL parsing, host validation, request creation, execution, unexpected status, response read, stream, and close errors in `attachment/service.go`.
- [x] 2.2 Add resolved attachment redirect URL context to missing `Location`, redirect parsing, host validation, and redirected download failures in `attachment/service.go`.
- [x] 2.3 Avoid changing caller validation errors that occur before any URL is accessed, such as missing writer, missing attachment link, or missing transport client.

## 3. Verification

- [x] 3.1 Update export tests to assert representative export start, redirect download, Cloud polling, and Cloud result download failures include the full attempted URL.
- [x] 3.2 Update attachment tests to assert representative direct download, redirect download, and unexpected status failures include the full attempted URL.
- [x] 3.3 Add or update tests proving existing sentinel export errors remain matchable with `errors.Is` after URL context is added.
- [x] 3.4 Run the focused package tests covering `export` and `attachment` behavior.
