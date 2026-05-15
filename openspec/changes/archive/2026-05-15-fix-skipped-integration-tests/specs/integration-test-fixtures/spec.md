## ADDED Requirements

### Requirement: Integration fixture data is deterministic
The integration test suite SHALL treat missing scenario fixture data as a test failure after the suite has been explicitly invoked with the integration build tag and required Confluence credentials.

#### Scenario: Required integration environment is absent
- **WHEN** an integration test is run without a required environment variable for credentials or fixture IDs
- **THEN** the test fails with a message naming the missing environment variable

#### Scenario: Required scenario data is absent
- **WHEN** an integration test reaches a scenario whose configured fixture exists but lacks data required for that scenario
- **THEN** the test fails with an actionable message describing the missing fixture data
- **AND** the test does not call `t.Skip`, `t.Skipf`, or equivalent skip behavior for that scenario

### Requirement: Attachment download uses a page with attachments
The integration test suite SHALL use a configured page fixture that contains at least one attachment when testing attachment download behavior.

#### Scenario: Dedicated attachment page is configured
- **WHEN** `CONFLUENCE_ATTACHMENT_PAGE_ID` is set
- **THEN** attachment download integration tests use that page ID to discover downloadable attachments

#### Scenario: Dedicated attachment page is not configured
- **WHEN** `CONFLUENCE_ATTACHMENT_PAGE_ID` is not set
- **THEN** attachment download integration tests use `CONFLUENCE_PAGE_ID` as the fallback page ID

#### Scenario: Attachment page has no attachments
- **WHEN** the page selected for attachment download discovery has no attachments
- **THEN** the test fails with a message identifying the page ID and instructing maintainers to configure `CONFLUENCE_ATTACHMENT_PAGE_ID` with a page that has at least one attachment

### Requirement: Integration fixture documentation is complete
The project SHALL document every environment variable needed to run the integration suite without scenario-level skips.

#### Scenario: Maintainer configures integration tests
- **WHEN** a maintainer reads the integration test setup documentation
- **THEN** the documentation lists `CONFLUENCE_ATTACHMENT_PAGE_ID`
- **AND** it states that the page must contain at least one attachment
- **AND** it states that the value may match `CONFLUENCE_PAGE_ID` only when that page satisfies the attachment requirement
