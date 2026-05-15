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
		return fmt.Errorf("download redirect from %q: %w", responseURL(resp, c.client.BaseURL), ErrMissingLocation)
	}

	parsedBase, err := url.Parse(c.client.BaseURL)
	if err != nil {
		return fmt.Errorf("parse baseURL %q: %w", c.client.BaseURL, err)
	}

	downloadURL, err := parsedBase.Parse(location)
	if err != nil {
		return fmt.Errorf("parse Location header %q from %q: %w", location, responseURL(resp, c.client.BaseURL), err)
	}

	err = c.ensureRedirectHostAllowed(downloadURL.Host)
	if err != nil {
		return fmt.Errorf("validate redirect URL %q: %w", downloadURL.String(), err)
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
		return fmt.Errorf("create download request for %q: %w", downloadURL, err)
	}

	pdfResp, err := c.client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute download request for %q: %w", downloadURL, err)
	}
	defer func() {
		if closeErr := pdfResp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close download response body for %q: %w", downloadURL, closeErr)
		}
	}()

	if pdfResp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(pdfResp.Body)
		if readErr != nil {
			return fmt.Errorf("read download error response from %q: %w", downloadURL, readErr)
		}
		return fmt.Errorf("unexpected download status code %d for %q: %s", pdfResp.StatusCode, downloadURL, string(body))
	}

	_, err = io.Copy(writer, pdfResp.Body)
	if err != nil {
		return fmt.Errorf("stream pdf response body from %q: %w", downloadURL, err)
	}

	return nil
}

func responseURL(resp *http.Response, fallback string) string {
	if resp != nil && resp.Request != nil && resp.Request.URL != nil {
		return resp.Request.URL.String()
	}
	return fallback
}
