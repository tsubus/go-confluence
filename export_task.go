package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProgressResponse struct {
	Progress               int    `json:"progress"`
	State                  string `json:"state"`
	Result                 string `json:"result"`
	EstimatedTimeRemaining int    `json:"estimatedTimeRemaining"`
	TimeElapsed            int    `json:"timeElapsed"`
}

func (c *Client) handleOKResponse(ctx context.Context, resp *http.Response, writer io.Writer) error {
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

func (c *Client) pollTaskProgress(ctx context.Context, taskID string) (string, error) {
	pollURL := fmt.Sprintf("%s/api/v2/pdfexporttask/progress/%s", c.baseURL, url.PathEscape(taskID))

	pollCtx := ctx
	var cancel context.CancelFunc
	if c.pollTimeout > 0 {
		pollCtx, cancel = context.WithTimeout(ctx, c.pollTimeout)
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

func (c *Client) fetchProgress(ctx context.Context, pollURL string) (ProgressResponse, error) {
	req, err := c.newRequest(ctx, pollURL)
	if err != nil {
		return ProgressResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
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

func (c *Client) waitForNextPoll(ctx context.Context) error {
	ticker := time.NewTimer(c.pollInterval)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while polling: %w", ctx.Err())
	case <-ticker.C:
		return nil
	}
}
