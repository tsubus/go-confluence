//go:build integration

package confluence

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestIntegration_ExportPage performs a live export against a real Confluence
// instance. It is skipped unless the required environment variables are set.
//
// Required environment variables:
//   - CONFLUENCE_URL      e.g. http://localhost:8090
//   - CONFLUENCE_USERNAME
//   - CONFLUENCE_PASSWORD
//   - CONFLUENCE_PAGE_ID
func TestIntegration_ExportPage(t *testing.T) {
	baseURL := os.Getenv("CONFLUENCE_URL")
	username := os.Getenv("CONFLUENCE_USERNAME")
	password := os.Getenv("CONFLUENCE_PASSWORD")
	pageID := os.Getenv("CONFLUENCE_PAGE_ID")

	if baseURL == "" || username == "" || password == "" || pageID == "" {
		t.Skip("Skipping integration test: set CONFLUENCE_URL, CONFLUENCE_USERNAME, CONFLUENCE_PASSWORD and CONFLUENCE_PAGE_ID")
	}

	client, err := NewClient(baseURL, WithBasicAuth(username, password))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pdfBytes, err := client.ExportPage(ctx, pageID)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)
	require.True(t, bytes.HasPrefix(pdfBytes, []byte("%PDF")), "ExportPage() returned non-PDF content: %q", pdfBytes[:min(len(pdfBytes), 20)])

	t.Logf("Successfully exported page %s (%d bytes)", pageID, len(pdfBytes))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
