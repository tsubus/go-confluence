package confluence

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/umats/go-confluence/attachment"
	"github.com/umats/go-confluence/internal/transport"
)

// BuildQueryForTest exposes transport.BuildQuery for external tests.
func BuildQueryForTest(params any) (url.Values, error) {
	values, err := transport.BuildQuery(params)
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	return values, nil
}

// NewURLForTest exposes v2 client URL construction for external tests.
func NewURLForTest(client *Client, path string) (string, error) {
	v2 := transport.NewV2Client(newTransportClient(client))
	url, err := transport.NewURLForTest(v2, path)
	if err != nil {
		return "", fmt.Errorf("build url: %w", err)
	}
	return url, nil
}

// DoJSONForTest exposes v2 client JSON request flow for external tests.
func DoJSONForTest(
	ctx context.Context,
	client *Client,
	method string,
	path string,
	query url.Values,
	request any,
	response any,
) error {
	v2 := transport.NewV2Client(newTransportClient(client))
	err := v2.DoJSON(ctx, method, path, query, request, response)
	if err != nil {
		return fmt.Errorf("do json request: %w", err)
	}
	return nil
}

// DecodeResponseForTest exposes v2 client response decoding for external tests.
func DecodeResponseForTest(client *Client, response *http.Response, out any) (transport.APIError, error) {
	v2 := transport.NewV2Client(newTransportClient(client))
	apiErr := transport.APIError{}
	err := transport.DecodeResponseForTest(v2, response, out)
	if err == nil {
		return apiErr, nil
	}
	if errors.As(err, &apiErr) {
		return apiErr, fmt.Errorf("decode response: %w", err)
	}
	return apiErr, fmt.Errorf("decode response: %w", err)
}

// ResetAttachmentDownloadByURLDeprecationWarningForTest resets the once-only deprecation warning for tests.
func ResetAttachmentDownloadByURLDeprecationWarningForTest() {
	attachment.ResetDeprecatedDownloadByURLWarningForTest()
}

// NewAttachmentDownloadByURLDeprecationWarningForTest exposes warning construction for external tests.
func NewAttachmentDownloadByURLDeprecationWarningForTest() interface {
	ResetForTest()
	Warn()
} {
	return attachment.NewDownloadByURLDeprecationWarningForTest()
}
