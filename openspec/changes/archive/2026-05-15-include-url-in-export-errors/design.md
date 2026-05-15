## Context

The client has two URL-accessing flows relevant to this change:

- Page PDF export in `export`: build an export URL, request it, handle Server/Data Center redirects, handle Cloud HTML task responses, poll task progress, then download the resulting PDF URL.
- Attachment download in `attachment`: fetch attachment metadata or accept a direct download URL, validate redirect/download hosts, follow download redirects, and stream the file.

Today many errors say what failed but not which URL was being accessed. That makes production failures hard to debug when redirects, Cloud task result URLs, proxy behavior, tenant base paths, or attachment links differ between environments.

## Goals / Non-Goals

**Goals:**
- Add URL context to every error produced while accessing, downloading, or polling a page export or attachment download URL.
- Include the resolved absolute URL when the code has one, not only the raw path or `Location` header.
- Preserve existing public method signatures and sentinel error compatibility with `errors.Is`.
- Keep implementation local to existing export and attachment flows.

**Non-Goals:**
- Introduce a new exported error type or structured logging API.
- Change redirect policy, authentication, polling cadence, timeout behavior, or response parsing.
- Redact or transform URLs beyond using the exact URL already attempted by the client.
- Broaden this behavior to all REST v2 JSON APIs.

## Decisions

1. **Wrap errors at the call site with `url %q` context.**
   - Use `fmt.Errorf("<operation> %q: %w", attemptedURL, err)` for wrapped errors.
   - Rationale: it preserves `errors.Is`/`errors.As` and follows the current code style.
   - Alternative considered: add a new URL-aware error type. Rejected because callers asked for diagnostic message context, and a new type would expand public API surface without a current need.

2. **Use resolved absolute URLs in messages.**
   - For export start requests, use `exportURL` from `Client.ExportURL(pageID)`.
   - For redirect downloads, resolve `Location` against the configured base URL before validation/download and use `downloadURL.String()` in downstream errors.
   - For Cloud polling, use the full `pollURL` constructed from `BaseURL` and `taskID`.
   - For attachment direct downloads, use the provided `downloadURL`; for redirects, use the resolved redirect URL.
   - Alternative considered: report only the raw `Location` header. Rejected because relative redirects are common and the full attempted URL is more useful for debugging.

3. **Do not include URLs in input validation errors before URL access.**
   - Examples: missing page ID, nil writer, nil transport client, empty attachment download link.
   - Rationale: these are caller validation failures rather than errors from accessing or downloading a URL.

4. **Keep sentinel errors wrapped, not replaced.**
   - `ErrMissingLocation`, `ErrTaskIDNotFound`, `ErrTaskResultEmpty`, and `ErrTaskFailed` must remain discoverable with `errors.Is` after URL context is added where applicable.
   - Alternative considered: concatenate sentinel error text into a new error string. Rejected because it would break existing sentinel checks.

## Risks / Trade-offs

- [Risk] URLs can contain sensitive query parameters in some deployments → Mitigation: this change follows the explicit debugging requirement to include the full URL; do not add partial redaction that would make diagnostics incomplete.
- [Risk] Tests that assert exact error strings may need updates → Mitigation: use substring and `errors.Is` assertions for behavior rather than exact full messages.
- [Risk] Double-wrapping could make messages noisy in redirect paths → Mitigation: wrap once at the boundary where the attempted URL is known; avoid wrapping the same URL repeatedly in helper recursion.
