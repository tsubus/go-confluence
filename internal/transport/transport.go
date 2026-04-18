package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Client holds shared transport configuration for subpackages.
type Client struct {
	BaseURL                  string
	HTTPClient               *http.Client
	Username                 string
	Password                 string
	PollInterval             time.Duration
	PollTimeout              time.Duration
	AllowedRedirectHosts     map[string]struct{}
	AllowCrossHostContentURL bool
}

// NewRequest builds an HTTP request with context and optional basic auth.
func (c *Client) NewRequest(ctx context.Context, method, requestURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
	return req, nil
}

// NewV2URL builds a REST v2 URL from the configured base URL.
func (c *Client) NewV2URL(p string) (string, error) {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("parse baseURL: %w", err)
	}
	base.Path = strings.TrimSuffix(base.Path, "/")
	if !strings.HasSuffix(base.Path, "/wiki") {
		base.Path = path.Join(base.Path, "wiki")
	}
	base.Path = path.Join(base.Path, "api", "v2", p)
	return base.String(), nil
}

// ExportURL builds the export endpoint for a given page.
func (c *Client) ExportURL(pageID string) (string, error) {
	exportURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("parse baseURL: %w", err)
	}
	exportURL.Path = path.Join(exportURL.Path, "spaces/flyingpdf/pdfpageexport.action")

	q := exportURL.Query()
	q.Set("pageId", pageID)
	exportURL.RawQuery = q.Encode()

	return exportURL.String(), nil
}
