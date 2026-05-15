## 1. Fixture Configuration

- [x] 1.1 Extend `integrationEnv` with `CONFLUENCE_ATTACHMENT_PAGE_ID` fallback handling
- [x] 1.2 Update integration test header comments with the attachment page fixture
- [x] 1.3 Add CI environment mapping for `CONFLUENCE_ATTACHMENT_PAGE_ID`
- [x] 1.4 Update README integration setup with attachment page guidance

## 2. Deterministic Scenario Coverage

- [x] 2.1 Replace attachment download skip path with an assertion failure
- [x] 2.2 Use the attachment page fixture when listing attachments for download
- [x] 2.3 Ensure failure message identifies the page ID and remediation variable
- [x] 2.4 Search integration tests for remaining data-dependent skip calls and remove any scenario-level skips

## 3. Verification

- [x] 3.1 Run targeted unit tests for integration fixture parsing if helper tests are added
- [x] 3.2 Run `go test -tags=integration -run TestIntegration_AttachmentDownload` with configured Confluence fixtures
- [x] 3.3 Confirm the selected attachment page ID has at least one attachment or update the fixture to a valid page ID
- [x] 3.4 Run `openspec status --change fix-skipped-integration-tests` and confirm apply readiness
