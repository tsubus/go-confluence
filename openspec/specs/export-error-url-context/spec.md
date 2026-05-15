## Requirements

### Requirement: Page export errors include attempted URL
The system SHALL include the full URL it attempted to access in errors produced while requesting, handling, polling, or downloading a page export.

#### Scenario: Export start request fails
- **WHEN** a page export request cannot be created or executed for a resolved export URL
- **THEN** the returned error message includes that full resolved export URL

#### Scenario: Export response has an unexpected status
- **WHEN** the page export endpoint returns an unexpected status code
- **THEN** the returned error message includes the full export URL that returned the response

#### Scenario: Export redirect download fails
- **WHEN** a Server/Data Center page export redirects to a PDF URL and that download fails
- **THEN** the returned error message includes the full resolved PDF download URL

#### Scenario: Cloud export polling fails
- **WHEN** a Cloud page export starts a task and polling the task progress endpoint fails
- **THEN** the returned error message includes the full task progress URL that was polled

#### Scenario: Cloud export result download fails
- **WHEN** a Cloud page export task returns a result URL and downloading that URL fails
- **THEN** the returned error message includes the full result download URL

### Requirement: Attachment download errors include attempted URL
The system SHALL include the full URL it attempted to access in errors produced while validating, requesting, redirecting, or streaming an attachment download.

#### Scenario: Direct attachment download fails
- **WHEN** an attachment direct download URL cannot be validated, requested, or streamed successfully
- **THEN** the returned error message includes that full direct download URL

#### Scenario: Attachment redirect download fails
- **WHEN** an attachment download response redirects to another URL and the redirected download fails
- **THEN** the returned error message includes the full resolved redirected download URL

#### Scenario: Attachment download response has an unexpected status
- **WHEN** an attachment download URL returns an unexpected status code
- **THEN** the returned error message includes the full URL that returned the response

### Requirement: Existing error classification is preserved
The system SHALL preserve existing sentinel error behavior while adding URL context to export and attachment download failures.

#### Scenario: Sentinel export error is wrapped with URL context
- **WHEN** a failure corresponds to an existing sentinel export error
- **THEN** callers can still match the sentinel with `errors.Is`
- **AND** the returned error message includes the relevant attempted URL when the failure occurred during URL access or download handling
