package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// APIError represents a non-2xx API response.
type APIError struct {
	StatusCode int
	Body       string
}

// Error implements the error interface.
func (e APIError) Error() string {
	return fmt.Sprintf("unexpected status code %d: %s", e.StatusCode, e.Body)
}

// V2Client provides helper methods for REST v2 APIs.
type V2Client struct {
	client *Client
}

// NewV2Client creates a V2Client from a transport Client.
func NewV2Client(client *Client) *V2Client {
	return &V2Client{client: client}
}

// Client exposes the underlying transport configuration.
func (v *V2Client) Client() *Client {
	return v.client
}

func (v *V2Client) newURL(p string) (string, error) {
	base, err := url.Parse(v.client.BaseURL)
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

// NewURLForTest exposes URL construction for external tests.
func NewURLForTest(v2 *V2Client, path string) (string, error) {
	return v2.newURL(path)
}

// DecodeResponseForTest exposes response decoding for external tests.
func DecodeResponseForTest(v2 *V2Client, response *http.Response, out any) error {
	return v2.decodeResponse(response, out)
}

func (v *V2Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if v.client.Username != "" {
		req.SetBasicAuth(v.client.Username, v.client.Password)
	}
	return req, nil
}

func (v *V2Client) do(req *http.Request) (*http.Response, error) {
	resp, err := v.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	if resp == nil {
		return nil, errors.New("request returned nil response")
	}
	return resp, nil
}

func (v *V2Client) decodeResponse(resp *http.Response, out any) error {
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}
	if out == nil {
		return nil
	}
	if len(body) == 0 {
		return nil
	}
	err = json.Unmarshal(body, out)
	if err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// DoJSON executes an HTTP request with JSON request/response handling.
func (v *V2Client) DoJSON(ctx context.Context, method, path string, query url.Values, request any, response any) error {
	endpoint, err := v.newURL(path)
	if err != nil {
		return err
	}
	if len(query) > 0 {
		endpoint = endpoint + "?" + query.Encode()
	}

	var body io.Reader
	if request != nil {
		payload, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			return fmt.Errorf("encode request: %w", marshalErr)
		}
		body = bytes.NewReader(payload)
	}

	req, err := v.newRequest(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	if request != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := v.do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return v.decodeResponse(resp, response)
}
