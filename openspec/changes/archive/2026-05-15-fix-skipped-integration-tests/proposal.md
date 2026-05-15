## Why

The integration suite currently skips coverage when the configured external Confluence fixture lacks required data, such as `CONFLUENCE_PAGE_ID=98380` having no attachments for the attachment download path. Skipped live tests hide regressions in API coverage that the integration suite is intended to catch.

## What Changes

- Replace data-dependent skip behavior in integration tests with explicit fixture requirements and failing diagnostics.
- Allow tests that need specialized fixtures, especially attachment download, to use a dedicated page ID when the general page fixture does not satisfy the scenario.
- Add preflight validation for integration environment fixtures so missing pages, spaces, labels, attachments, or page-attached files fail early with actionable messages instead of skipping during scenario execution.
- Document any additional required integration environment variable and selection guidance for a page that contains at least one downloadable attachment.
- No public client API changes.

## Capabilities

### New Capabilities
- `integration-test-fixtures`: Requirements for deterministic Confluence integration test fixtures and scenario coverage without condition-based skips.

### Modified Capabilities
- `export-error-url-context`: No requirement changes.

## Impact

- Affected code: `client_integration_test.go` and any integration-test documentation or CI configuration that describes required Confluence fixture variables.
- Affected systems: external Confluence test instance used by the `integration` build tag.
- APIs and dependencies: no production API or dependency changes expected.
