package confluence

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// ExportPage exports a single Confluence page as a PDF.
// It returns the raw PDF bytes or an error if the export fails.
func (c *Client) ExportPage(ctx context.Context, pageID string) ([]byte, error) {
	var buffer bytes.Buffer
	err := c.ExportPageTo(ctx, pageID, &buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ExportPageTo exports a single Confluence page as a PDF into the provided writer.
func (c *Client) ExportPageTo(ctx context.Context, pageID string, writer io.Writer) (err error) {
	if strings.TrimSpace(pageID) == "" {
		return errors.New("pageID is required")
	}
	if writer == nil {
		return errors.New("writer is required")
	}

	exportURL, err := c.exportURL(pageID)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, exportURL)
	if err != nil {
		return err
	}

	req.Header.Set("X-Atlassian-Token", "no-check")
	req.Header.Set("Accept", "application/json,text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute export request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close export response body: %w", closeErr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusFound:
		return c.downloadFromRedirect(ctx, resp, writer)
	case http.StatusOK:
		return c.handleOKResponse(ctx, resp, writer)
	default:
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("read export error response: %w", readErr)
		}
		return fmt.Errorf("unexpected export status code %d: %s", resp.StatusCode, string(body))
	}
}

func (c *Client) exportURL(pageID string) (string, error) {
	exportURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse baseURL: %w", err)
	}
	exportURL.Path = path.Join(exportURL.Path, "spaces/flyingpdf/pdfpageexport.action")

	q := exportURL.Query()
	q.Set("pageId", pageID)
	exportURL.RawQuery = q.Encode()

	return exportURL.String(), nil
}
