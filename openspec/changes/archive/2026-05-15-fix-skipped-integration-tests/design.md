## Context

The repository has a live Confluence integration suite in `client_integration_test.go` behind the `integration` build tag. Most tests exercise endpoints using shared fixture variables from the environment. `TestIntegration_AttachmentDownload` currently derives a page attachment by listing attachments on `CONFLUENCE_PAGE_ID`; when that page has no attachments it calls `t.Skipf`, which makes a green integration run possible without exercising attachment download. The CI workflow and README currently expose one page fixture and one attachment fixture, but no explicit page-with-attachment fixture.

## Goals / Non-Goals

**Goals:**
- Make integration runs deterministic: once required Confluence credentials and fixture IDs are configured, scenario-specific missing fixture data must fail with actionable diagnostics rather than skip.
- Preserve the suite-level ability to not run integration tests unless the `integration` build tag and required environment are supplied.
- Support a dedicated page ID for attachment-download scenarios so the general page fixture can remain focused on page API coverage.
- Keep production client code and public APIs unchanged.

**Non-Goals:**
- Provision or mutate Confluence test data automatically.
- Replace live integration tests with mocked HTTP tests.
- Validate every possible Confluence content relationship beyond the fixtures required by existing integration scenarios.

## Decisions

1. Introduce an optional-but-recommended `CONFLUENCE_ATTACHMENT_PAGE_ID` fixture for page-attachment scenarios.
   - If set, attachment download uses it instead of `CONFLUENCE_PAGE_ID`.
   - If unset, the code may fall back to `CONFLUENCE_PAGE_ID` for local compatibility, but an empty attachment list must fail with a message instructing maintainers to set `CONFLUENCE_ATTACHMENT_PAGE_ID` to a page containing at least one attachment.
   - Alternative considered: require all consumers to change `CONFLUENCE_PAGE_ID` to a page with attachments. Rejected because the general page fixture may need to represent a stable page for other page scenarios, and changing it can disturb unrelated tests.

2. Replace data-dependent `t.Skipf` calls with assertions/failures after environment preconditions are met.
   - Missing environment variables remain hard failures via `requireIntegrationEnv` because invoking integration tests is an explicit opt-in.
   - Scenario data gaps, such as a page with zero attachments, are reported as fixture misconfiguration.
   - Alternative considered: dynamically search the entire Confluence instance for a matching page. Rejected because search can be slow, permission-dependent, and less deterministic than an explicit fixture ID.

3. Add a fixture preflight helper only for relationships whose absence would otherwise skip or invalidate a scenario.
   - The helper should reuse existing client calls and return concrete IDs needed by tests, avoiding duplicate setup logic inside individual tests.
   - It should not create abstractions for unrelated endpoint smoke tests that already fail correctly on missing top-level IDs.
   - Alternative considered: a comprehensive global preflight for every endpoint. Rejected because it adds latency and can duplicate the actual test assertions without improving the current skipped-scenario problem.

4. Update documentation and CI environment mappings alongside test code.
   - README and the GitHub Actions integration-test environment should mention `CONFLUENCE_ATTACHMENT_PAGE_ID`.
   - The docs should state that the value may be the same as `CONFLUENCE_PAGE_ID` only if that page has at least one attachment.

## Risks / Trade-offs

- Fixture data may be deleted or permissions may change in the external Confluence instance → tests fail loudly with fixture remediation instructions rather than silently skipping coverage.
- Adding a dedicated fixture variable increases setup burden → mitigated by allowing it to match the existing page ID when that page already has attachments.
- Live integration tests remain dependent on external Confluence state → this change improves determinism of configured scenarios but does not remove the external dependency.
