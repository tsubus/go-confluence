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
//   - CONFLUENCE_LABEL_ID
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

func TestIntegration_PageGetAttachments(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	attachments, err := client.Page().GetAttachments(ctx, pageID, &PageAttachmentsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, attachments)
}

func TestIntegration_PageGetLabels(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	labels, err := client.Page().GetLabels(ctx, pageID, &PageLabelsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, labels)
}

func TestIntegration_PageGetOperations(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	ops, err := client.Page().GetOperations(ctx, pageID)
	require.NoError(t, err)
	require.NotNil(t, ops)
}

func TestIntegration_PageGetVersions(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	versions, err := client.Page().GetVersions(ctx, pageID, &PageVersionsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, versions)
}

func TestIntegration_PageGetVersionDetails(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	version, err := client.Page().GetVersionDetails(ctx, pageID, 1)
	require.NoError(t, err)
	require.NotNil(t, version)
}

func TestIntegration_PageGetAncestors(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	ancestors, err := client.Page().GetAncestors(ctx, pageID, &PageAncestorsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, ancestors)
}

func TestIntegration_PageGetDescendants(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	descendants, err := client.Page().GetDescendants(ctx, pageID, nil)
	require.NoError(t, err)
	require.NotNil(t, descendants)
}

func TestIntegration_PageGetChildren(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	children, err := client.Page().GetChildren(ctx, pageID, &PageChildrenParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, children)
}

func TestIntegration_PageGetDirectChildren(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	children, err := client.Page().GetDirectChildren(ctx, pageID, &PageChildrenParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, children)
}

func TestIntegration_PageGetFooterComments(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	comments, err := client.Page().GetFooterComments(ctx, pageID, &PageCommentsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, comments)
}

func TestIntegration_PageGetInlineComments(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	comments, err := client.Page().GetInlineComments(ctx, pageID, &PageCommentsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, comments)
}

func TestIntegration_PageGetContentProperties(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	props, err := client.Page().GetContentProperties(ctx, pageID, &ContentPropertiesParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, props)
}

func TestIntegration_PageGetLikeCount(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	likes, err := client.Page().GetLikeCount(ctx, pageID)
	require.NoError(t, err)
	require.NotNil(t, likes)
}

func TestIntegration_PageGetLikeUsers(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 1
	users, err := client.Page().GetLikeUsers(ctx, pageID, &PageLikeUsersParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, users)
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

func TestIntegration_SpaceGetPages(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	pages, err := client.Space().GetPages(ctx, env.spaceID, &SpacePagesParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, pages)
}

func TestIntegration_SpaceGetBlogPosts(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	posts, err := client.Space().GetBlogPosts(ctx, env.spaceID, &SpaceBlogPostsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, posts)
}

func TestIntegration_SpaceGetLabels(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	labels, err := client.Space().GetLabels(ctx, env.spaceID, &SpaceLabelsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, labels)
}

func TestIntegration_SpaceGetContentLabels(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	labels, err := client.Space().GetContentLabels(ctx, env.spaceID, &SpaceLabelsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, labels)
}

func TestIntegration_SpaceGetOperations(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ops, err := client.Space().GetOperations(ctx, env.spaceID)
	require.NoError(t, err)
	require.NotNil(t, ops)
}

func TestIntegration_SpaceGetPermissions(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	perms, err := client.Space().GetPermissions(ctx, env.spaceID, &SpacePermissionsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, perms)
}

func TestIntegration_SpaceGetProperties(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	props, err := client.Space().GetProperties(ctx, env.spaceID, &SpacePropertiesParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, props)
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

func TestIntegration_AttachmentGetLabels(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	labels, err := client.Attachment().GetLabels(ctx, env.attachmentID, &AttachmentLabelsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, labels)
}

func TestIntegration_AttachmentGetOperations(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ops, err := client.Attachment().GetOperations(ctx, env.attachmentID)
	require.NoError(t, err)
	require.NotNil(t, ops)
}

func TestIntegration_AttachmentGetVersions(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	versions, err := client.Attachment().GetVersions(ctx, env.attachmentID, &AttachmentVersionsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, versions)
}

func TestIntegration_AttachmentGetVersionDetails(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	version, err := client.Attachment().GetVersionDetails(ctx, env.attachmentID, 1)
	require.NoError(t, err)
	require.NotNil(t, version)
}

func TestIntegration_AttachmentGetComments(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	comments, err := client.Attachment().GetComments(ctx, env.attachmentID, &AttachmentCommentsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, comments)
}

func TestIntegration_AttachmentGetContentProperties(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	props, err := client.Attachment().GetContentProperties(ctx, env.attachmentID, &ContentPropertiesParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, props)
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

func TestIntegration_LabelGetAttachments(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	attachments, err := client.Label().GetAttachments(ctx, env.labelID, &LabelAttachmentsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, attachments)
}

func TestIntegration_LabelGetBlogPosts(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	limit := 1
	posts, err := client.Label().GetBlogPosts(ctx, env.labelID, &LabelBlogPostsParams{Limit: &limit})
	require.NoError(t, err)
	require.NotNil(t, posts)
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

func TestIntegration_AttachmentDownload(t *testing.T) {
	env := requireIntegrationEnv(t)

	client, err := NewClient(env.baseURL, WithBasicAuth(env.username, env.password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Find attachments on the configured page.
	pageID, err := strconv.ParseInt(env.pageID, 10, 64)
	require.NoError(t, err)

	limit := 10
	attachments, err := client.Page().GetAttachments(
		ctx,
		pageID,
		&PageAttachmentsParams{Limit: &limit},
	)
	require.NoError(t, err)
	require.NotNil(t, attachments)
	if len(attachments.Results) == 0 {
		t.Skipf("page %s has no attachments; skipping download test", env.pageID)
	}

	first := attachments.Results[0]
	require.NotNil(t, first.Id, "attachment has no id")

	// Download the first attachment.
	var buf bytes.Buffer
	err = client.Attachment().Download(ctx, *first.Id, nil, &buf)
	require.NoError(t, err)
	require.NotZero(t, buf.Len(), "downloaded attachment is empty")

	t.Logf("Downloaded attachment %s (%d bytes)", *first.Id, buf.Len())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
