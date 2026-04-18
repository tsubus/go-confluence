package export

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Helper) downloadFromRedirect(ctx context.Context, resp *http.Response, writer io.Writer) error {
	location := resp.Header.Get("Location")
	if location == "" {
		return ErrMissingLocation
	}

	parsedBase, err := url.Parse(c.client.BaseURL)
	if err != nil {
		return fmt.Errorf("parse baseURL: %w", err)
	}

	downloadURL, err := parsedBase.Parse(location)
	if err != nil {
		return fmt.Errorf("parse Location header %q: %w", location, err)
	}

	err = c.ensureRedirectHostAllowed(downloadURL.Host)
	if err != nil {
		return err
	}

	return c.downloadPDF(ctx, downloadURL.String(), writer)
}

func (c *Helper) ensureRedirectHostAllowed(host string) error {
	if host == "" {
		return errors.New("redirect host is empty")
	}
	if len(c.client.AllowedRedirectHosts) == 0 {
		return errors.New("allowed redirect host list is empty")
	}
	_, ok := c.client.AllowedRedirectHosts[host]
	if !ok {
		return fmt.Errorf("redirect host %q is not allowed", host)
	}
	return nil
}

func (c *Helper) downloadPDF(ctx context.Context, downloadURL string, writer io.Writer) (err error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("create download request: %w", err)
	}

	pdfResp, err := c.client.HTTPClient.Do(req)
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
