// Package confluence provides a minimal client for exporting Confluence
// pages as PDFs. It supports both Confluence Server/Data Center and
// Confluence Cloud installations.
//
// Server/Data Center: the export action returns a 302 redirect to the PDF.
// Confluence Cloud: the export action starts a background task; the client
// polls the task progress endpoint and downloads the PDF when ready.
package confluence

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	defaultTimeout      = 30 * time.Second
	defaultPollInterval = 3 * time.Second
	minTaskIDMatches    = 2
	progressComplete    = 100
)

var (
	taskIDRegex = regexp.MustCompile(`<meta[^>]+name="ajs-taskId"[^>]+content="([^"]+)"`)
	// ErrMissingLocation indicates the export response lacked a Location header.
	ErrMissingLocation = errors.New("export response missing Location header")
	// ErrTaskFailed indicates the Confluence export task failed.
	ErrTaskFailed = errors.New("pdf export task failed")
	// ErrTaskResultEmpty indicates the task finished without a result URL.
	ErrTaskResultEmpty = errors.New("task completed but result URL is empty")
	// ErrTaskIDNotFound indicates the Cloud export HTML lacked a task ID.
	ErrTaskIDNotFound = errors.New("taskId meta tag not found")
)

// Client exports Confluence pages as PDFs.
type Client struct {
	baseURL      string
	httpClient   *http.Client
	username     string
	password     string
	pollInterval time.Duration
	pollTimeout  time.Duration
	requireHTTPS bool
}

// HTTPDoer abstracts [http.Client] for testing and future REST APIs.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Option configures a Client.
type Option func(*Client) error

// WithBasicAuth sets the username and password (or API token) for Basic Auth.
func WithBasicAuth(username, password string) Option {
	return func(c *Client) error {
		c.username = username
		c.password = password
		return nil
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithPollInterval overrides the default poll interval for Cloud background tasks.
func WithPollInterval(d time.Duration) Option {
	return func(c *Client) error {
		if d <= 0 {
			return errors.New("poll interval must be positive")
		}
		c.pollInterval = d
		return nil
	}
}

// WithPollTimeout sets the maximum time spent polling a Cloud export task.
func WithPollTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d <= 0 {
			return errors.New("poll timeout must be positive")
		}
		c.pollTimeout = d
		return nil
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d <= 0 {
			return errors.New("timeout must be positive")
		}
		c.httpClient.Timeout = d
		return nil
	}
}

// WithRequireHTTPS enforces https scheme for baseURL.
func WithRequireHTTPS() Option {
	return func(c *Client) error {
		c.requireHTTPS = true
		return nil
	}
}

// NewClient creates a new Confluence PDF export client.
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("baseURL is required")
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("baseURL must have a scheme and host")
	}

	// Ensure no trailing slash so path joining works predictably.
	parsed.Path = strings.TrimSuffix(parsed.Path, "/")

	c := &Client{
		baseURL: parsed.String(),
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			// We handle the 302 redirect manually to capture the Location header.
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		pollInterval: defaultPollInterval,
	}

	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	if c.requireHTTPS && parsed.Scheme != "https" {
		return nil, errors.New("baseURL must use https")
	}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, requestURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	return req, nil
}
