//go:build integration

package confluence

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestIntegration_ExportPage performs a live export against a real Confluence
// instance. It fails fast unless the required environment variables are set.
//
// Required environment variables:
//   - CONFLUENCE_URL e.g. http://localhost:8090
//   - CONFLUENCE_USERNAME
//   - CONFLUENCE_PASSWORD
//   - CONFLUENCE_PAGE_ID
//   - CONFLUENCE_SPACE_ID
//   - CONFLUENCE_ATTACHMENT_ID
func TestIntegration_ExportPage(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pdfBytes, err := client.Export().Page(ctx, env.pageID)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)
	require.True(t, bytes.HasPrefix(pdfBytes, []byte("%PDF")), "Export().Page() returned non-PDF content: %q", pdfBytes[:min(len(pdfBytes), 20)])

	t.Logf("Successfully exported page %s (%d bytes)", env.pageID, len(pdfBytes))
}

func TestIntegration_PageGet(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	page, err := client.Page().Get(ctx, pageID, nil)
	require.NoError(t, err)
	require.NotNil(t, page)
}

func TestIntegration_PageList(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	pages, err := client.Page().List(ctx, &PageListParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, pages)
}

func TestIntegration_SpaceGet(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	space, err := client.Space().Get(ctx, env.spaceID, nil)
	require.NoError(t, err)
	require.NotNil(t, space)
}

func TestIntegration_SpaceList(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	spaces, err := client.Space().List(ctx, &SpaceListParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, spaces)
}

func TestIntegration_AttachmentGet(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	attachment, err := client.Attachment().Get(ctx, env.attachmentID, nil)
	require.NoError(t, err)
	require.NotNil(t, attachment)
}

func TestIntegration_AttachmentList(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	attachments, err := client.Attachment().List(ctx, &AttachmentListParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, attachments)
}

func TestIntegration_LabelGetPages(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pages, err := client.Label().GetPages(ctx, env.labelID, &LabelPagesParams{})
	require.NoError(t, err)
	require.NotNil(t, pages)
}

func TestIntegration_LabelList(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	labels, err := client.Label().List(ctx, &LabelListParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, labels)
}

type integrationEnv struct {
	baseURL      string
	username     string
	password     string
	pageID       string
	spaceID      int64
	attachmentID string
	labelID      int64
}

func requireIntegrationEnv(t *testing.T) integrationEnv {
	t.Helper()

	baseURL := os.Getenv("CONFLUENCE_URL")
	username := os.Getenv("CONFLUENCE_USERNAME")
	password := os.Getenv("CONFLUENCE_PASSWORD")
	pageID := os.Getenv("CONFLUENCE_PAGE_ID")
	spaceIDValue := os.Getenv("CONFLUENCE_SPACE_ID")
	attachmentID := os.Getenv("CONFLUENCE_ATTACHMENT_ID")
	labelIDValue := os.Getenv("CONFLUENCE_LABEL_ID")

	require.NotEmpty(t, baseURL, "CONFLUENCE_URL must be set")
	require.NotEmpty(t, username, "CONFLUENCE_USERNAME must be set")
	require.NotEmpty(t, password, "CONFLUENCE_PASSWORD must be set")
	require.NotEmpty(t, pageID, "CONFLUENCE_PAGE_ID must be set")
	require.NotEmpty(t, spaceIDValue, "CONFLUENCE_SPACE_ID must be set")
	require.NotEmpty(t, attachmentID, "CONFLUENCE_ATTACHMENT_ID must be set")
	require.NotEmpty(t, labelIDValue, "CONFLUENCE_LABEL_ID must be set")

	spaceID, err := strconv.ParseInt(spaceIDValue, 10, 64)
	require.NoError(t, err)

	labelID, err := strconv.ParseInt(labelIDValue, 10, 64)
	require.NoError(t, err)

	return integrationEnv{
		baseURL:      baseURL,
		username:     username,
		password:     password,
		pageID:       pageID,
		spaceID:      spaceID,
		attachmentID: attachmentID,
		labelID:      labelID,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
