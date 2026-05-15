## MODIFIED Requirements

### Requirement: Attachment download by ID uses the supported API endpoint
The system SHALL download attachments by ID through Confluence's supported attachment download endpoint and SHALL NOT fall back to attachment metadata download links when that endpoint fails.

#### Scenario: Supported API download succeeds
- **WHEN** attachment metadata is resolved and the supported attachment download endpoint returns the attachment content successfully
- **THEN** `AttachmentService.Download` streams that content to the caller

#### Scenario: Supported API download fails before streaming
- **WHEN** attachment metadata is resolved and the supported attachment download endpoint fails before any bytes are written
- **THEN** `AttachmentService.Download` returns that failure directly
- **AND** it does not attempt the attachment metadata `downloadLink` as a fallback

#### Scenario: Supported API download fails after partial streaming
- **WHEN** the supported attachment download endpoint writes bytes to the caller and then fails
- **THEN** `AttachmentService.Download` returns that streaming failure directly
- **AND** it does not attempt any alternate download path

### Requirement: Direct attachment URL download is deprecated
The system SHALL keep direct attachment URL download available temporarily for caller-supplied URLs, but SHALL mark it deprecated and warn callers when it is used.

#### Scenario: Caller uses direct URL download helper
- **WHEN** `AttachmentService.DownloadByURL` is called
- **THEN** the call still attempts to stream the provided URL using the existing direct-download behavior
- **AND** the system emits a deprecation warning indicating that `DownloadByURL` is deprecated

#### Scenario: ID-based attachment download does not emit direct URL deprecation warning
- **WHEN** `AttachmentService.Download` is called for attachment retrieval by ID
- **THEN** the system does not emit the `DownloadByURL` deprecation warning as part of that code path
