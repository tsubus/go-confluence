package export

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/umats/go-confluence/internal/transport"
)

const (
	minTaskIDMatches = 2
	progressComplete = 100
)

var (
	taskIDRegex        = regexp.MustCompile(`<meta[^>]+name="ajs-taskId"[^>]+content="([^"]+)"`)
	ErrMissingLocation = errors.New("export response missing Location header")
	ErrTaskFailed      = errors.New("pdf export task failed")
	ErrTaskResultEmpty = errors.New("task completed but result URL is empty")
	ErrTaskIDNotFound  = errors.New("taskId meta tag not found")
)

// Helper holds shared export dependencies.
type Helper struct {
	client *transport.Client
}

// NewHelper wraps a transport client for export helpers.
func NewHelper(client *transport.Client) *Helper {
	return &Helper{client: client}
}

// DownloadFromRedirect exposes redirect download handling for tests.
func (c *Helper) DownloadFromRedirect(ctx context.Context, resp *http.Response, writer io.Writer) error {
	return c.downloadFromRedirect(ctx, resp, writer)
}

// DownloadPDF exposes PDF download for tests.
func (c *Helper) DownloadPDF(ctx context.Context, downloadURL string, writer io.Writer) error {
	return c.downloadPDF(ctx, downloadURL, writer)
}

// HandleOKResponse exposes export response handling for tests.
func (c *Helper) HandleOKResponse(ctx context.Context, resp *http.Response, writer io.Writer) error {
	return c.handleOKResponse(ctx, resp, writer)
}

// FetchProgress exposes task progress fetching for tests.
func (c *Helper) FetchProgress(ctx context.Context, pollURL string) (ProgressResponse, error) {
	return c.fetchProgress(ctx, pollURL)
}

// WaitForNextPoll exposes polling wait for tests.
func (c *Helper) WaitForNextPoll(ctx context.Context) error {
	return c.waitForNextPoll(ctx)
}

// PollTaskProgress exposes task polling for tests.
func (c *Helper) PollTaskProgress(ctx context.Context, taskID string) (string, error) {
	return c.pollTaskProgress(ctx, taskID)
}

type ProgressResponse struct {
	Progress               int    `json:"progress"`
	State                  string `json:"state"`
	Result                 string `json:"result"`
	EstimatedTimeRemaining int    `json:"estimatedTimeRemaining"`
	TimeElapsed            int    `json:"timeElapsed"`
}

func (c *Helper) handleOKResponse(ctx context.Context, resp *http.Response, writer io.Writer) error {
	var body []byte
	var err error

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/pdf") {
		_, err = io.Copy(writer, resp.Body)
		if err != nil {
			return fmt.Errorf("stream pdf response body: %w", err)
		}
		return nil
	}

	// Assume HTML indicating a Cloud background task.
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read export response body: %w", err)
	}

	taskID, err := extractTaskID(string(body))
	if err != nil {
		return fmt.Errorf("extract task ID from HTML: %w", err)
	}

	downloadURL, err := c.pollTaskProgress(ctx, taskID)
	if err != nil {
		return fmt.Errorf("poll task progress: %w", err)
	}

	return c.downloadPDF(ctx, downloadURL, writer)
}

func extractTaskID(html string) (string, error) {
	matches := taskIDRegex.FindStringSubmatch(html)
	if len(matches) < minTaskIDMatches {
		return "", ErrTaskIDNotFound
	}
	return matches[1], nil
}

// ExtractTaskIDForTest exposes task ID extraction for fuzz tests.
func ExtractTaskIDForTest(html string) (string, error) {
	return extractTaskID(html)
}

func (c *Helper) pollTaskProgress(ctx context.Context, taskID string) (string, error) {
	pollURL := fmt.Sprintf("%s/api/v2/pdfexporttask/progress/%s", c.client.BaseURL, url.PathEscape(taskID))

	pollCtx := ctx
	var cancel context.CancelFunc
	if c.client.PollTimeout > 0 {
		pollCtx, cancel = context.WithTimeout(ctx, c.client.PollTimeout)
		defer cancel()
	}

	for attempt := 1; ; attempt++ {
		pr, err := c.fetchProgress(pollCtx, pollURL)
		if err != nil {
			return "", err
		}

		result, done, err := evaluateProgress(pr)
		if err != nil {
			return "", err
		}
		if done {
			return result, nil
		}

		waitErr := c.waitForNextPoll(pollCtx)
		if waitErr != nil {
			return "", fmt.Errorf("poll attempt %d: %w", attempt, waitErr)
		}
	}
}

func (c *Helper) fetchProgress(ctx context.Context, pollURL string) (ProgressResponse, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, pollURL, nil)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("create poll request: %w", err)
	}

	resp, err := c.client.HTTPClient.Do(req)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("execute poll request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("read poll response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ProgressResponse{}, fmt.Errorf("unexpected poll status code %d: %s", resp.StatusCode, string(body))
	}

	var pr ProgressResponse
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("decode poll response: %w", err)
	}

	return pr, nil
}

func evaluateProgress(pr ProgressResponse) (string, bool, error) {
	if pr.Progress >= progressComplete {
		if pr.State == "SUCCEEDED" || pr.State == "UPLOADED_TO_S3" {
			if pr.Result == "" {
				return "", true, ErrTaskResultEmpty
			}
			return pr.Result, true, nil
		}
		if pr.State == "FAILED" {
			return "", true, ErrTaskFailed
		}
		return "", true, fmt.Errorf("pdf export task finished with unexpected state: %s", pr.State)
	}

	if pr.State == "FAILED" {
		return "", true, ErrTaskFailed
	}

	return "", false, nil
}

func (c *Helper) waitForNextPoll(ctx context.Context) error {
	ticker := time.NewTimer(c.client.PollInterval)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while polling: %w", ctx.Err())
	case <-ticker.C:
		return nil
	}
}
