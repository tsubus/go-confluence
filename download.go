package confluence

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) downloadFromRedirect(ctx context.Context, resp *http.Response, writer io.Writer) error {
	location := resp.Header.Get("Location")
	if location == "" {
		return ErrMissingLocation
	}

	parsedBase, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("parse baseURL: %w", err)
	}

	downloadURL, err := parsedBase.Parse(location)
	if err != nil {
		return fmt.Errorf("parse Location header %q: %w", location, err)
	}

	return c.downloadPDF(ctx, downloadURL.String(), writer)
}

func (c *Client) downloadPDF(ctx context.Context, downloadURL string, writer io.Writer) (err error) {
	req, err := c.newRequest(ctx, downloadURL)
	if err != nil {
		return err
	}

	pdfResp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute download request: %w", err)
	}
	defer func() {
		if closeErr := pdfResp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close download response body: %w", closeErr)
		}
	}()

	if pdfResp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(pdfResp.Body)
		if readErr != nil {
			return fmt.Errorf("read download error response: %w", readErr)
		}
		return fmt.Errorf("unexpected download status code %d: %s", pdfResp.StatusCode, string(body))
	}

	_, err = io.Copy(writer, pdfResp.Body)
	if err != nil {
		return fmt.Errorf("stream pdf response body: %w", err)
	}

	return nil
}
