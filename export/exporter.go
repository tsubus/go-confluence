package export

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/umats/go-confluence/internal/transport"
)

type exporter struct {
	client *transport.Client
	helper *Helper
}

func newExporter(client *transport.Client) *exporter {
	return &exporter{client: client, helper: NewHelper(client)}
}

func (e *exporter) Page(ctx context.Context, pageID string) ([]byte, error) {
	var buffer bytes.Buffer
	err := e.PageTo(ctx, pageID, &buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (e *exporter) PageTo(ctx context.Context, pageID string, writer io.Writer) (err error) {
	if strings.TrimSpace(pageID) == "" {
		return errors.New("pageID is required")
	}
	if writer == nil {
		return errors.New("writer is required")
	}

	exportURL, err := e.client.ExportURL(pageID)
	if err != nil {
		return fmt.Errorf("build export URL: %w", err)
	}

	req, err := e.client.NewRequest(ctx, http.MethodGet, exportURL, nil)
	if err != nil {
		return fmt.Errorf("create export request for %q: %w", exportURL, err)
	}

	req.Header.Set("X-Atlassian-Token", "no-check")
	req.Header.Set("Accept", "application/json,text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := e.client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute export request for %q: %w", exportURL, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close export response body for %q: %w", exportURL, closeErr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusFound:
		return e.helper.downloadFromRedirect(ctx, resp, writer)
	case http.StatusOK:
		return e.helper.handleOKResponse(ctx, resp, writer)
	default:
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("read export error response from %q: %w", exportURL, readErr)
		}
		return fmt.Errorf("unexpected export status code %d for %q: %s", resp.StatusCode, exportURL, string(body))
	}
}
